package peer_material_repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PendingMaterialReceiveRequestRepositorySql struct {
	db *sql.DB
}

func MakePendingMaterialReceiveRequestRepositorySql(
	iDb *sql.DB,
) PendingMaterialReceiveRequestRepositorySql {
	return PendingMaterialReceiveRequestRepositorySql{
		db: iDb,
	}
}

func (r PendingMaterialReceiveRequestRepositorySql) CreatePendingReceiveMaterialRequest(
	iContext context.Context,
	iUser models.User,
	iToBeReceivedMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iSenderPublicKey string,
	iTransferTime models.CustomTime,
	iIsOutbound bool,
) (models.PendingMaterialReceiveRequest, error) {
	tx, err := r.db.BeginTx(iContext, nil)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	requestResponse := tx.QueryRowContext(
		iContext,
		`INSERT INTO "pending_receive_material_request" (
			main_node_id,
			transfer_time,
			sender_public_key,
			owner_id,
			is_outbound
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING id
	`, *iToBeReceivedMaterial.Id, iTransferTime, iSenderPublicKey, iUser.Id, iIsOutbound)

	var pendingRequestId int
	err = requestResponse.Scan(&pendingRequestId)
	if err != nil {
		tx.Rollback()
		return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create pending request %w", err)
	}

	if len(iRelatedMaterials) > 0 {
		queryValString := make([]string, 0, len(iRelatedMaterials))
		queryVal := make([]interface{}, 0, len(iRelatedMaterials)*2)
		count := 0
		for _, material := range iRelatedMaterials {
			queryValString = append(queryValString, fmt.Sprintf("($%d, %d)", count, count+1))
			queryVal = append(queryVal, *&material.Id)
			queryVal = append(queryVal, pendingRequestId)
			count += 2
		}
		query := fmt.Sprintf(`
			INSERT INTO "node_from_peer" (
				node_id,
				pending_receive_material_request_id
			) VALUES %s
		`, strings.Join(queryValString, ","))

		_, err = tx.ExecContext(
			iContext,
			query,
			queryVal...,
		)
		if err != nil {
			tx.Rollback()
			return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create signature options %w", err)
		}
	}

	if len(iOptions) > 0 {
		count := 1
		queryValString := make([]string, 0, len(iOptions))
		queryVal := make([]interface{}, 0, len(iOptions)*3)
		for _, option := range iOptions {
			queryValString = append(queryValString, fmt.Sprintf("($%d, $%d, $%d)", count, count+1, count+2))
			queryVal = append(queryVal, option.Signature)
			queryVal = append(queryVal, option.Id)
			queryVal = append(queryVal, pendingRequestId)
			count += 3
		}

		query := fmt.Sprintf(`
			INSERT INTO "signature_option" (
				signature,
				new_node_id,
				pending_request_id
			) VALUES %s
		`, strings.Join(queryValString, ","))

		_, err = tx.ExecContext(
			iContext,
			query,
			queryVal...,
		)
		if err != nil {
			tx.Rollback()
			return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create signature options %w", err)
		}
	}

	tx.Commit()

	ret := models.MakePendingMaterialReceiveRequest(
		pendingRequestId,
		iUser.Id,
		iToBeReceivedMaterial,
		iRelatedMaterials,
		iOptions,
		iSenderPublicKey,
		iTransferTime,
	)

	return ret, err
}

func (r PendingMaterialReceiveRequestRepositorySql) FetchPendingReceiveMaterialRequestsByUserId(
	iContext context.Context,
	iUserId int,
	iIsOutbound bool,
) ([]SimplifiedPendingReceiveMaterialRequest, error) {
	response, err := r.db.Query(`
		SELECT 
			request.id,
			request.transfer_time,
			request.sender_public_key,
			request.main_node_id
		FROM "pending_receive_material_request" request
		WHERE request.owner_id=$1 AND request.is_outbound=$2
	`, iUserId, iIsOutbound)

	if err != nil {
		return []SimplifiedPendingReceiveMaterialRequest{}, err
	}
	defer response.Close()

	type RequestInfo struct {
		requestId       int
		senderPublicKey string
		transferTime    models.CustomTime
		mainMaterialId  int
	}

	fetchedRequests := map[int]RequestInfo{}

	for response.Next() {
		request := RequestInfo{}
		err := response.Scan(
			&request.requestId,
			&request.transferTime,
			&request.senderPublicKey,
			&request.mainMaterialId,
		)
		if err != nil {
			return []SimplifiedPendingReceiveMaterialRequest{}, err
		}

		fetchedRequests[request.requestId] = request
	}

	if len(fetchedRequests) == 0 {
		return []SimplifiedPendingReceiveMaterialRequest{}, nil
	}

	relatedMaterialIdsMap := map[int][]int{}
	{
		/// fetch nodes from peer
		argStr := []string{}
		arg := []interface{}{}
		count := 1
		for id := range fetchedRequests {
			argStr = append(argStr, fmt.Sprintf("(pending_receive_material_request_id=$%d)", count))
			arg = append(arg, id)
			count++
		}
		query := `
			SELECT 
				node_id,
				pending_receive_material_request_id
			FROM "node_from_peer"
			WHERE 
		`
		query += strings.Join(argStr, " OR ")
		response, err := r.db.Query(query, arg...)
		if err != nil {
			return []SimplifiedPendingReceiveMaterialRequest{}, nil
		}
		defer response.Close()
		for response.Next() {
			nodeId := 0
			pendingReceiveMaterialRequestId := 0
			err := response.Scan(
				&nodeId,
				&pendingReceiveMaterialRequestId,
			)
			if err != nil {
				return []SimplifiedPendingReceiveMaterialRequest{}, nil
			}

			if _, ok := relatedMaterialIdsMap[pendingReceiveMaterialRequestId]; ok {
				relatedMaterialIdsMap[pendingReceiveMaterialRequestId] = []int{}
			}
			relatedMaterialIdsMap[pendingReceiveMaterialRequestId] = append(relatedMaterialIdsMap[pendingReceiveMaterialRequestId], nodeId)
		}
	}

	signatureOptionsMap := map[int][]models.SignatureOption{}
	{
		/// fetch nodes from peer
		argStr := []string{}
		arg := []interface{}{}
		count := 1
		for id := range fetchedRequests {
			argStr = append(argStr, fmt.Sprintf("(pending_request_id=$%d)", count))
			arg = append(arg, id)
			count++
		}
		query := `
			SELECT 
				signature,
				new_node_id,
				pending_request_id
			FROM "signature_option"
			WHERE 
		`
		query += strings.Join(argStr, " OR ")
		response, err := r.db.Query(query, arg...)
		if err != nil {
			return []SimplifiedPendingReceiveMaterialRequest{}, nil
		}
		defer response.Close()
		for response.Next() {
			signature := ""
			newNodeId := ""
			pendingRequestId := 0
			err := response.Scan(
				&signature,
				&newNodeId,
				&pendingRequestId,
			)
			if err != nil {
				return []SimplifiedPendingReceiveMaterialRequest{}, nil
			}

			if _, ok := signatureOptionsMap[pendingRequestId]; ok {
				signatureOptionsMap[pendingRequestId] = []models.SignatureOption{}
			}
			signatureOption := models.MakeSignatureOption(newNodeId, signature)
			signatureOptionsMap[pendingRequestId] = append(signatureOptionsMap[pendingRequestId], signatureOption)
		}
	}

	ret := []SimplifiedPendingReceiveMaterialRequest{}
	for requestId := range fetchedRequests {
		relatedMaterialIds := []int{}
		if fetchedRelatedMaterialIds, ok := relatedMaterialIdsMap[requestId]; ok {
			relatedMaterialIds = fetchedRelatedMaterialIds
		}

		signatureOptions := []models.SignatureOption{}
		if fetchedSignatureOptions, ok := signatureOptionsMap[requestId]; ok {
			signatureOptions = fetchedSignatureOptions
		}

		request := makeSimplifiedPendingReceiveMaterialRequest(
			requestId,
			iUserId,
			fetchedRequests[requestId].mainMaterialId,
			relatedMaterialIds,
			signatureOptions,
			fetchedRequests[requestId].senderPublicKey,
			fetchedRequests[requestId].transferTime,
		)

		ret = append(ret, request)
	}

	return ret, nil
}
