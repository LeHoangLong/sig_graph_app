package peer_material_repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type MaterialReceiveRequestRepositorySql struct {
	db *sql.DB
}

func MakeMaterialReceiveRequestRepositorySql(
	iDb *sql.DB,
) MaterialReceiveRequestRepositorySql {
	return MaterialReceiveRequestRepositorySql{
		db: iDb,
	}
}

func (r MaterialReceiveRequestRepositorySql) createReceiveMaterialRequest(
	iContext context.Context,
	tx *sql.Tx,
	iToBeReceivedMaterialId models.NodeId,
	iTransferTime models.CustomTime,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iStatus models.MaterialReceiveRequestStatus,
) (models.MaterialReceiveRequestId, error) {
	response := tx.QueryRowContext(
		iContext,
		`INSERT INTO "receive_material_request" (
			main_node_id,
			transfer_time,
			status_id
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING id`,
		iToBeReceivedMaterialId,
		iTransferTime,
		iStatus,
	)

	pendingReceiveMaterialRequestId := models.MaterialReceiveRequestId(0)
	err := response.Scan(&pendingReceiveMaterialRequestId)
	if err != nil {
		return models.MaterialReceiveRequestId(0), err
	}

	if len(iRelatedMaterials) > 0 {
		queryValString := make([]string, 0, len(iRelatedMaterials))
		queryVal := make([]interface{}, 0, len(iRelatedMaterials)*2)
		count := 0
		for _, material := range iRelatedMaterials {
			queryValString = append(queryValString, fmt.Sprintf("($%d, %d)", count, count+1))
			queryVal = append(queryVal, *material.Id)
			queryVal = append(queryVal, pendingReceiveMaterialRequestId)
			count += 2
		}
		query := fmt.Sprintf(`
			INSERT INTO "node_from_peer" (
				node_id,
				receive_material_request_id
			) VALUES %s
		`, strings.Join(queryValString, ","))

		_, err = tx.ExecContext(
			iContext,
			query,
			queryVal...,
		)
		if err != nil {
			return models.MaterialReceiveRequestId(0), fmt.Errorf("cannot create signature options %w", err)
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
			queryVal = append(queryVal, pendingReceiveMaterialRequestId)
			count += 3
		}

		query := fmt.Sprintf(`
			INSERT INTO "signature_option" (
				signature,
				new_node_id,
				request_id
			) VALUES %s
		`, strings.Join(queryValString, ","))

		_, err = tx.ExecContext(
			iContext,
			query,
			queryVal...,
		)
		if err != nil {
			return models.MaterialReceiveRequestId(0), fmt.Errorf("cannot create signature options %w", err)
		}
	}

	return pendingReceiveMaterialRequestId, nil
}

func (r MaterialReceiveRequestRepositorySql) CreateOutboundReceiveMaterialRequest(
	iContext context.Context,
	iRecipientPeerId models.PeerId,
	iSenderUserId models.UserId,
	iSenderPublicKeyId models.PublicKeyId,
	iToBeReceivedMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iTransferTime models.CustomTime,
	iStatus models.MaterialReceiveRequestStatus,
) (models.OutboundMaterialReceiveRequest, error) {
	tx, err := r.db.BeginTx(iContext, nil)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, fmt.Errorf("cannot start transaction %w", err)
	}
	pendingReceiveMaterialRequestId, err := r.createReceiveMaterialRequest(
		iContext,
		tx,
		*iToBeReceivedMaterial.Id,
		iTransferTime,
		iRelatedMaterials,
		iOptions,
		iStatus,
	)
	if err != nil {
		tx.Rollback()
		return models.OutboundMaterialReceiveRequest{}, err
	}

	_, err = tx.QueryContext(
		iContext,
		`INSERT INTO "outbound_receive_material_request" (
			request_id,
			owner_id,
			owner_public_key_id,
			recipient_peer_id
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`, pendingReceiveMaterialRequestId, iSenderUserId, iSenderPublicKeyId, iRecipientPeerId)

	if err != nil {
		tx.Rollback()
		return models.OutboundMaterialReceiveRequest{}, fmt.Errorf("cannot create pending request %w", err)
	}

	tx.Commit()

	pendingRequest := models.MakeMaterialReceiveRequest(
		pendingReceiveMaterialRequestId,
		iToBeReceivedMaterial,
		iRelatedMaterials,
		iOptions,
		iTransferTime,
		iStatus,
	)

	outboundRequest := models.MakeOutboundMaterialReceiveRequest(
		pendingRequest,
		iRecipientPeerId,
		iSenderUserId,
		iSenderPublicKeyId,
	)

	return outboundRequest, err
}

func (r MaterialReceiveRequestRepositorySql) CreateInboundReceiveMaterialRequest(
	iContext context.Context,
	iRecipientUserId models.UserId,
	iSenderPublicKeyId models.PublicKeyId,
	iToBeReceivedMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iTransferTime models.CustomTime,
	iStatus models.MaterialReceiveRequestStatus,
) (models.InboundMaterialReceiveRequest, error) {
	tx, err := r.db.BeginTx(iContext, nil)
	if err != nil {
		return models.InboundMaterialReceiveRequest{}, fmt.Errorf("cannot start transaction %w", err)
	}
	pendingReceiveMaterialRequestId, err := r.createReceiveMaterialRequest(
		iContext,
		tx,
		*iToBeReceivedMaterial.Id,
		iTransferTime,
		iRelatedMaterials,
		iOptions,
		iStatus,
	)
	if err != nil {
		tx.Rollback()
		return models.InboundMaterialReceiveRequest{}, err
	}

	_, err = tx.QueryContext(
		iContext,
		`INSERT INTO "inbound_receive_material_request" (
			request_id,
			owner_id,
			sender_public_key_id
		) VALUES (
			$1,
			$2,
			$3
		)
	`, pendingReceiveMaterialRequestId, iRecipientUserId, iSenderPublicKeyId)
	if err != nil {
		tx.Rollback()
		return models.InboundMaterialReceiveRequest{}, err
	}

	tx.Commit()

	pendingRequest := models.MakeMaterialReceiveRequest(
		pendingReceiveMaterialRequestId,
		iToBeReceivedMaterial,
		iRelatedMaterials,
		iOptions,
		iTransferTime,
		iStatus,
	)

	inboundRequest := models.MakeInboundMaterialReceiveRequest(
		pendingRequest,
		iRecipientUserId,
		iSenderPublicKeyId,
	)
	return inboundRequest, nil
}

func (r MaterialReceiveRequestRepositorySql) FetchInboundReceiveMaterialRequestsByUserId(
	iContext context.Context,
	iUserId models.UserId,
	iStatus []models.MaterialReceiveRequestStatus,
) ([]SimplifiedInboundReceiveMaterialRequest, error) {
	query := `
		SELECT 
			parent.id,
			parent.transfer_time,
			request.sender_public_key_id,
			parent.main_node_id,
			parent.status_id
		FROM "inbound_receive_material_request" request
		INNER JOIN "receive_material_request" parent
			ON parent.id = request.request_id
			AND request.owner_id=$1
	`
	arg := []interface{}{iUserId}
	if len(iStatus) > 0 {
		count := 2
		argString := []string{}
		for _, status := range iStatus {
			argString = append(argString, fmt.Sprintf("(parent.status_id=$%d)", count))
			arg = append(arg, status)
			count += 1
		}
		query = query + " WHERE " + strings.Join(argString, " OR ")
	}
	response, err := r.db.Query(query, arg...)

	if err != nil {
		return []SimplifiedInboundReceiveMaterialRequest{}, err
	}
	defer response.Close()

	type RequestInfo struct {
		requestId         models.MaterialReceiveRequestId
		senderPublicKeyId models.PublicKeyId
		transferTime      models.CustomTime
		mainMaterialId    models.NodeId
		status            models.MaterialReceiveRequestStatus
	}

	fetchedRequests := map[models.MaterialReceiveRequestId]RequestInfo{}

	for response.Next() {
		request := RequestInfo{}
		err := response.Scan(
			&request.requestId,
			&request.transferTime,
			&request.senderPublicKeyId,
			&request.mainMaterialId,
			&request.status,
		)
		if err != nil {
			return []SimplifiedInboundReceiveMaterialRequest{}, err
		}

		fetchedRequests[request.requestId] = request
	}

	if len(fetchedRequests) == 0 {
		return []SimplifiedInboundReceiveMaterialRequest{}, nil
	}

	relatedMaterialIdsMap := map[models.MaterialReceiveRequestId][]models.NodeId{}
	{
		/// fetch nodes from peer
		argStr := []string{}
		arg := []interface{}{}
		count := 1
		for id := range fetchedRequests {
			argStr = append(argStr, fmt.Sprintf("(receive_material_request_id=$%d)", count))
			arg = append(arg, id)
			count++
		}
		query := `
			SELECT 
				node_id,
				receive_material_request_id
			FROM "node_from_peer"
			WHERE 
		`
		query += strings.Join(argStr, " OR ")
		response, err := r.db.Query(query, arg...)
		if err != nil {
			return []SimplifiedInboundReceiveMaterialRequest{}, nil
		}
		defer response.Close()
		for response.Next() {
			nodeId := models.NodeId(0)
			pendingReceiveMaterialRequestId := models.MaterialReceiveRequestId(0)
			err := response.Scan(
				&nodeId,
				&pendingReceiveMaterialRequestId,
			)
			if err != nil {
				return []SimplifiedInboundReceiveMaterialRequest{}, nil
			}

			if _, ok := relatedMaterialIdsMap[pendingReceiveMaterialRequestId]; ok {
				relatedMaterialIdsMap[pendingReceiveMaterialRequestId] = []models.NodeId{}
			}
			relatedMaterialIdsMap[pendingReceiveMaterialRequestId] = append(relatedMaterialIdsMap[pendingReceiveMaterialRequestId], nodeId)
		}
	}

	signatureOptionsMap := map[models.MaterialReceiveRequestId][]models.SignatureOption{}
	{
		/// fetch nodes from peer
		argStr := []string{}
		arg := []interface{}{}
		count := 1
		for id := range fetchedRequests {
			argStr = append(argStr, fmt.Sprintf("(request_id=$%d)", count))
			arg = append(arg, id)
			count++
		}
		query := `
			SELECT 
				signature,
				new_node_id,
				request_id
			FROM "signature_option"
			WHERE 
		`
		query += strings.Join(argStr, " OR ")
		response, err := r.db.Query(query, arg...)
		if err != nil {
			return []SimplifiedInboundReceiveMaterialRequest{}, nil
		}
		defer response.Close()
		for response.Next() {
			signature := ""
			newNodeId := ""
			pendingRequestId := models.MaterialReceiveRequestId(0)
			err := response.Scan(
				&signature,
				&newNodeId,
				&pendingRequestId,
			)
			if err != nil {
				return []SimplifiedInboundReceiveMaterialRequest{}, nil
			}

			if _, ok := signatureOptionsMap[pendingRequestId]; ok {
				signatureOptionsMap[pendingRequestId] = []models.SignatureOption{}
			}
			signatureOption := models.MakeSignatureOption(newNodeId, signature)
			signatureOptionsMap[pendingRequestId] = append(signatureOptionsMap[pendingRequestId], signatureOption)
		}
	}

	ret := []SimplifiedInboundReceiveMaterialRequest{}
	for requestId := range fetchedRequests {
		relatedMaterialIds := []models.NodeId{}
		if fetchedRelatedMaterialIds, ok := relatedMaterialIdsMap[requestId]; ok {
			relatedMaterialIds = fetchedRelatedMaterialIds
		}

		signatureOptions := []models.SignatureOption{}
		if fetchedSignatureOptions, ok := signatureOptionsMap[requestId]; ok {
			signatureOptions = fetchedSignatureOptions
		}

		request := makeSimplifiedReceiveMaterialRequest(
			requestId,
			iUserId,
			fetchedRequests[requestId].mainMaterialId,
			relatedMaterialIds,
			signatureOptions,
			fetchedRequests[requestId].senderPublicKeyId,
			fetchedRequests[requestId].transferTime,
			fetchedRequests[requestId].status,
		)

		ret = append(ret, request)
	}

	return ret, nil
}
