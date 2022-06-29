package peer_material_repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type MaterialReceiveRequestRepositorySql struct {
	db                     *sql.DB
	peerProtocolRepository repositories.PeerProtocolRepositoryI
}

func MakeMaterialReceiveRequestRepositorySql(
	iDb *sql.DB,
	iPeerProtocolRepository repositories.PeerProtocolRepositoryI,
) MaterialReceiveRequestRepositorySql {
	return MaterialReceiveRequestRepositorySql{
		db:                     iDb,
		peerProtocolRepository: iPeerProtocolRepository,
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
	iSenderEndpoints []models.SenderEndpoint,
) (models.InboundMaterialReceiveRequest, error) {
	if len(iSenderEndpoints) == 0 {
		return models.InboundMaterialReceiveRequest{}, common.InvalidArgument
	}

	tx, err := r.db.BeginTx(iContext, nil)
	if err != nil {
		return models.InboundMaterialReceiveRequest{}, fmt.Errorf("cannot start transaction %w", err)
	}
	receiveMaterialRequestId, err := r.createReceiveMaterialRequest(
		iContext,
		tx,
		*iToBeReceivedMaterial.Id,
		iTransferTime,
		iRelatedMaterials,
		iOptions,
		iStatus,
	)
	fmt.Println("receiveMaterialRequestId")
	fmt.Println(receiveMaterialRequestId)
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
	`, receiveMaterialRequestId, iRecipientUserId, iSenderPublicKeyId)
	if err != nil {
		tx.Rollback()
		return models.InboundMaterialReceiveRequest{}, err
	}

	argString := []string{}
	arg := []interface{}{receiveMaterialRequestId}
	count := 2
	for i := range iSenderEndpoints {
		arg = append(arg, iSenderEndpoints[i].Url, iSenderEndpoints[i].Protocol.Id)
		argString = append(argString, fmt.Sprintf("($1, $%d, $%d)", count, count+1))
		count += 2
	}

	{
		query := `INSERT INTO "inbound_receive_material_request_sender_endpoint" (
			request_id,
			url,
			protocol_id
		) VALUES `
		query += strings.Join(argString, ",")
		_, err = r.db.QueryContext(
			iContext,
			query,
			arg...,
		)
		if err != nil {
			tx.Rollback()
			return models.InboundMaterialReceiveRequest{}, err
		}
	}

	tx.Commit()

	pendingRequest := models.MakeMaterialReceiveRequest(
		receiveMaterialRequestId,
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
		iSenderEndpoints,
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

	/// fetch sender endpoints
	senderEndpoints := map[models.MaterialReceiveRequestId][]models.SenderEndpoint{}
	{
		type TempSenderEndpoint struct {
			RequestId  models.MaterialReceiveRequestId
			Url        string
			ProtocolId models.PeerProtocolId
		}

		tempSenderEndpoints := map[models.MaterialReceiveRequestId][]TempSenderEndpoint{}

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
				request_id,
				url,
				protocol_id
			FROM "inbound_receive_material_request_sender_endpoint"
			WHERE 
		`
		query += strings.Join(argStr, " OR ")
		response, err := r.db.QueryContext(
			iContext,
			query,
			arg...,
		)
		if err != nil {
			return []SimplifiedInboundReceiveMaterialRequest{}, nil
		}
		defer response.Close()
		for response.Next() {
			requestId := models.MaterialReceiveRequestId(0)
			url := ""
			protocolId := models.PeerProtocolId(0)
			err = response.Scan(
				&requestId,
				&url,
				&protocolId,
			)
			if err != nil {
				return []SimplifiedInboundReceiveMaterialRequest{}, nil
			}
			tempSenderEndpoints[requestId] = append(tempSenderEndpoints[requestId], TempSenderEndpoint{
				RequestId:  requestId,
				Url:        url,
				ProtocolId: protocolId,
			})
		}

		peerProtocolIds := map[models.PeerProtocolId]bool{}
		for key := range tempSenderEndpoints {
			for i := range tempSenderEndpoints[key] {
				peerProtocolIds[tempSenderEndpoints[key][i].ProtocolId] = true
			}
		}

		peerProtocols, err := r.peerProtocolRepository.FetchPeerProtocolByIds(
			iContext,
			peerProtocolIds,
		)
		if err != nil {
			return []SimplifiedInboundReceiveMaterialRequest{}, nil
		}

		for requestId := range tempSenderEndpoints {
			for i := range tempSenderEndpoints[requestId] {
				protocolId := tempSenderEndpoints[requestId][i].ProtocolId
				url := tempSenderEndpoints[requestId][i].Url
				if peerProtocol, ok := peerProtocols[protocolId]; ok {
					senderEndpoints[requestId] = append(senderEndpoints[requestId],
						models.MakeSenderEndpoint(peerProtocol, url),
					)
				}
			}
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

		request := makeSimplifiedInboundReceiveMaterialRequest(
			requestId,
			iUserId,
			fetchedRequests[requestId].mainMaterialId,
			relatedMaterialIds,
			signatureOptions,
			fetchedRequests[requestId].senderPublicKeyId,
			fetchedRequests[requestId].transferTime,
			fetchedRequests[requestId].status,
			senderEndpoints[requestId],
		)

		ret = append(ret, request)
	}

	return ret, nil
}

func (r MaterialReceiveRequestRepositorySql) UpdateMaterialReceiveRequestStatus(
	iContext context.Context,
	iRequestId models.MaterialReceiveRequestId,
	iStatus models.MaterialReceiveRequestStatus,
) error {
	response := r.db.QueryRowContext(
		iContext,
		`
			UPDATE "receive_material_request" SET status_id=$1 WHERE id=$2 RETURNING id
		`, iStatus, iRequestId,
	)

	temp := models.MaterialReceiveRequestId(0)
	err := response.Scan(&temp)

	if err != nil {
		if err == sql.ErrNoRows {
			return common.NotFound
		}
		return err
	}

	return nil
}
