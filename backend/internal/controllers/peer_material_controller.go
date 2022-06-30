package controllers

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/repositories"
	peer_material_repositories "backend/internal/repositories/peer_material"
	"backend/internal/services"
	material_contract_service "backend/internal/services/material_contract"
	"backend/internal/services/node_contract"
	peer_material_services "backend/internal/services/peer_material"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type PeerMaterialController struct {
	userService                                     services.UserService
	materialFetchService                            material_contract_service.MaterialFetchServiceI
	materialTransferService                         material_contract_service.MaterialTransferServiceI
	idHasher                                        node_contract.IdHasherI
	pendingMaterialReceiveRequestRepositoryService  *peer_material_services.ReceiveMaterialRequestRepositoryService
	peerClientFactory                               peer_material_services.PeerMaterialClientServiceFactory
	userKeyRepository                               repositories.UserKeyRepositoryI
	optionGenerator                                 material_contract_service.SignatureOptionGenerator
	materialRepositoryService                       services.MaterialRepositoryService
	pendingMaterialReceiveRequestResponseRepository peer_material_repositories.MaterialReceiveRequestAcknowledgementRepositoryI
	peerRepository                                  repositories.PeerRepositoryI
	peerKeyRepository                               repositories.PeerKeyRepositoryI
	userEndpointRepository                          repositories.UserEndpointRepositoryI
	materialReceiveRequestAcknowledgementRepository peer_material_repositories.MaterialReceiveRequestAcknowledgementRepositoryI
}

func MakePeerMaterialController(
	iUserService services.UserService,
	iMaterialFetchService material_contract_service.MaterialFetchServiceI,
	iIdHasher node_contract.IdHasherI,
	iMaterialReceiveRequestRepositoryService *peer_material_services.ReceiveMaterialRequestRepositoryService,
	iPeerClientFactory peer_material_services.PeerMaterialClientServiceFactory,
	iUserKeyRepository repositories.UserKeyRepositoryI,
	iOptionGenerator material_contract_service.SignatureOptionGenerator,
	iMaterialRepositoryService services.MaterialRepositoryService,
	iMaterialReceiveRequestResponseRepository peer_material_repositories.MaterialReceiveRequestAcknowledgementRepositoryI,
	iPeerRepository repositories.PeerRepositoryI,
	iPeerKeyRepository repositories.PeerKeyRepositoryI,
	iMaterialTransferService material_contract_service.MaterialTransferServiceI,
	iUserEndpointRepository repositories.UserEndpointRepositoryI,
	iMaterialReceiveRequestAcknowledgementRepository peer_material_repositories.MaterialReceiveRequestAcknowledgementRepositoryI,
) *PeerMaterialController {
	return &PeerMaterialController{
		userService:          iUserService,
		materialFetchService: iMaterialFetchService,
		idHasher:             iIdHasher,
		pendingMaterialReceiveRequestRepositoryService: iMaterialReceiveRequestRepositoryService,
		peerClientFactory:         iPeerClientFactory,
		userKeyRepository:         iUserKeyRepository,
		optionGenerator:           iOptionGenerator,
		materialRepositoryService: iMaterialRepositoryService,
		pendingMaterialReceiveRequestResponseRepository: iMaterialReceiveRequestResponseRepository,
		peerRepository:          iPeerRepository,
		peerKeyRepository:       iPeerKeyRepository,
		materialTransferService: iMaterialTransferService,
		userEndpointRepository:  iUserEndpointRepository,
		materialReceiveRequestAcknowledgementRepository: iMaterialReceiveRequestAcknowledgementRepository,
	}
}

func (c *PeerMaterialController) Handle(
	iContext context.Context,
	iRequest peer_material_services.ReceiveMaterialRequestRequest,
) (peer_material_services.ReceiveMaterialRequestResponse, error) {
	user, err := c.userService.FindUserWithPublicKey(iRequest.RecipientPublicKey)
	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, err
	}

	toBeReceivedMaterial, err := c.materialFetchService.GetMaterialById(iRequest.MainNodeId)
	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, err
	}

	fetchedMaterials := []models.Material{toBeReceivedMaterial}
	for nodeId := range iRequest.Nodes {
		hashedId := c.idHasher.Hash(nodeId)
		found := false
		for _, fetchedMaterial := range fetchedMaterials {
			for childHashedId := range fetchedMaterial.NextNodeHashedIds {
				if hashedId == childHashedId {
					found = true
					break
				}
			}

			if !found {
				for parentHashedId := range fetchedMaterial.PreviousNodeHashedIds {
					if hashedId == parentHashedId {
						found = true
						break
					}
				}
			}
		}

		if found {
			relatedMaterial, err := c.materialFetchService.GetMaterialById(nodeId)
			if err != nil {
				return peer_material_services.ReceiveMaterialRequestResponse{}, err
			}
			fetchedMaterials = append(fetchedMaterials, relatedMaterial)
		}
	}

	materialOwnerKeyMap := map[string]bool{}
	for index := range fetchedMaterials {
		materialOwnerKeyMap[fetchedMaterials[index].OwnerPublicKey.Value] = true
	}

	materialOwnerKeys := []string{}
	for key := range materialOwnerKeyMap {
		materialOwnerKeys = append(materialOwnerKeys, key)
	}

	materialOwnerKeyModels, err := c.peerKeyRepository.CreateOrFetchPeerKeysByValue(
		iContext,
		user.Id,
		materialOwnerKeys,
	)
	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, err
	}

	materialOwnerKeyModelsIdMap := map[string]models.PublicKeyId{}
	for index := range materialOwnerKeyModels {
		materialOwnerKeyModelsIdMap[materialOwnerKeyModels[index].Value] = *materialOwnerKeyModels[index].Id
	}

	for index := range fetchedMaterials {
		id := materialOwnerKeyModelsIdMap[fetchedMaterials[index].OwnerPublicKey.Value]
		fetchedMaterials[index].OwnerPublicKey.Id = &id
	}

	namespace := services.UserIdToNamespace(user.Id)
	fetchedMaterials, err = c.materialRepositoryService.SaveMaterialsIgnoreId(
		iContext,
		namespace,
		fetchedMaterials,
	)
	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, err
	}

	unlinkedMaterials := map[models.NodeId]models.Material{}
	for i := range fetchedMaterials {
		unlinkedMaterials[*fetchedMaterials[i].Id] = fetchedMaterials[i]
	}
	linkedMaterials, err := c.materialRepositoryService.LinkMaterials(
		iContext,
		namespace,
		unlinkedMaterials,
	)
	fetchedMaterials = []models.Material{}
	for id := range linkedMaterials {
		fetchedMaterials = append(fetchedMaterials, linkedMaterials[id])
	}

	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, err
	}

	options := []models.SignatureOption{}
	for _, option := range iRequest.SignatureOptions {
		options = append(options, models.MakeSignatureOption(option.NodeId, option.Signature))
	}

	senderKeyId := materialOwnerKeyModelsIdMap[iRequest.SenderPublicKey]
	pendingRequest, err := c.pendingMaterialReceiveRequestRepositoryService.CreateInboundReceiveMaterialRequest(
		iContext,
		user.Id,
		senderKeyId,
		fetchedMaterials[0],
		fetchedMaterials[1:],
		options,
		models.CustomTime(iRequest.TransferTime),
		iRequest.SenderEndpoints,
	)

	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, nil
	}

	response, err := c.pendingMaterialReceiveRequestResponseRepository.SaveOutboundAcknowledgement(
		pendingRequest.Id,
		models.MaterialReceiveRequestResponseId(strconv.Itoa(int(pendingRequest.Id))),
	)

	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, nil
	}

	return peer_material_services.MakeReceiveMaterialRequestResponse(
		string(response.ResponseId),
		true,
	), nil
}

func (c *PeerMaterialController) SendRequest(
	iContext context.Context,
	iRecipientPublicKeyId models.PublicKeyId,
	iMainMaterialId models.NodeId,
	iRelatedMaterialId []models.NodeId,
) (models.OutboundMaterialReceiveRequest, error) {
	senderId, err := services.GetCurrentUserFromContext(iContext)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	sender, err := c.userService.GetUserById(iContext, senderId)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	peer, err := c.peerRepository.FetchPeerByKeyId(iContext, iRecipientPublicKeyId)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	recipientPublicKey := ""
	for _, key := range peer.PublicKey {
		if *key.Id == iRecipientPublicKeyId {
			recipientPublicKey = key.Value
		}
	}

	if recipientPublicKey == "" {
		return models.OutboundMaterialReceiveRequest{}, common.NotFound
	}

	transferTime := models.CustomTime(time.Now())
	key, err := c.userKeyRepository.FetchDefaultUserKeyPair(sender.Id)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	mainMaterial, relatedMaterials, err := c.materialRepositoryService.FetchMaterialsAndRelated(iContext, iMainMaterialId)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}
	filteredRelatedMaterials := []models.Material{}
	for _, material := range relatedMaterials {
		for _, id := range iRelatedMaterialId {
			if id == *material.Id {
				filteredRelatedMaterials = append(filteredRelatedMaterials, material)
				break
			}
		}
	}

	options, err := c.optionGenerator.GenerateSignatureOptionsForMaterialTransfer(
		iContext,
		10,
		mainMaterial,
		key,
	)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	request, err := c.pendingMaterialReceiveRequestRepositoryService.CreateOutboundReceiveMaterialRequest(
		iContext,
		peer.Id,
		senderId,
		*key.PublicKey.Id,
		mainMaterial,
		filteredRelatedMaterials,
		options,
		transferTime,
	)

	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	nodes := peer_material_services.GenerateNodesFromMaterials(append(filteredRelatedMaterials, mainMaterial))
	nodesMap := map[string]peer_material_services.Node{}
	for _, node := range nodes {
		nodesMap[node.Id] = node
	}

	serviceLayerOptions := make([]peer_material_services.SignatureOption, 0, len(options))
	for _, option := range options {
		serviceLayerOptions = append(serviceLayerOptions, peer_material_services.MakeNodeSignatureOption(
			option.Id,
			option.Signature,
		))
	}

	senderKeys, err := c.userKeyRepository.FetchUserKeyPairByUser(senderId)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	keyToUse := models.UserKeyPair{}
	keyFound := false
	for i := range senderKeys {
		if *senderKeys[i].Id == request.SenderPublicKeyId {
			keyToUse = senderKeys[i]
			keyFound = true
			break
		}
	}

	if !keyFound {
		return models.OutboundMaterialReceiveRequest{}, common.NotFound
	}

	userEndpoints, err := c.userEndpointRepository.FetchUserEndpointByUserId(
		iContext,
		senderId,
	)

	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	senderEndpoints := []peer_material_services.SenderEndpoint{}
	for i := range userEndpoints {
		senderEndpoint := peer_material_services.MakeSenderEndpoint(
			string(userEndpoints[i].Protocol.Name),
			userEndpoints[i].Protocol.MajorVersion,
			userEndpoints[i].Protocol.MinorVersion,
			userEndpoints[i].Url,
		)
		senderEndpoints = append(senderEndpoints, senderEndpoint)
	}

	peerRequest := peer_material_services.MakeReceiveMaterialRequestRequest(
		recipientPublicKey,
		mainMaterial.NodeId,
		nodesMap,
		time.Time(transferTime),
		keyToUse.PublicKey.Value,
		serviceLayerOptions,
		senderEndpoints,
	)

	endpoints, err := c.peerRepository.FetchPeerEndPoints(iContext, peer.Id)
	if err != nil {
		return models.OutboundMaterialReceiveRequest{}, err
	}

	requestReceived := false
	for _, endpoint := range endpoints {
		client, err := c.peerClientFactory.BuildPeerMaterialClientService(
			endpoint.Protocol.Name,
			endpoint.Protocol.MajorVersion,
			endpoint.Protocol.MinorVersion,
			endpoint.Url,
		)
		if err != nil {
			continue
		}

		response, err := client.SendReceiveMaterialRequest(iContext, peerRequest)
		if err != nil {
			return models.OutboundMaterialReceiveRequest{}, err
		}

		if !response.IsRequestAcknowledged {
			return models.OutboundMaterialReceiveRequest{}, fmt.Errorf("request rejected")
		}

		_, err = c.pendingMaterialReceiveRequestResponseRepository.SaveInboundAcknowledgement(
			request.Id,
			models.MaterialReceiveRequestResponseId(response.ResponseId),
		)
		if err != nil {
			return models.OutboundMaterialReceiveRequest{}, err
		}

		requestReceived = true
		break
	}

	if !requestReceived {
		return models.OutboundMaterialReceiveRequest{}, errors.New("could not send request")
	}

	return request, nil
}

func (c *PeerMaterialController) FetchReceivedMaterialReceiveRequests(
	iContext context.Context,
	iUserId models.UserId,
	iStatus []models.MaterialReceiveRequestStatus,
) ([]models.InboundMaterialReceiveRequest, error) {
	return c.pendingMaterialReceiveRequestRepositoryService.FetchInboundReceiveMaterialRequestsByUser(
		iContext,
		iUserId,
		iStatus,
	)
}

func (c *PeerMaterialController) AcceptPendingMaterialReceiveRequest(
	iContext context.Context,
	iUserId models.UserId,
	iRequestId models.MaterialReceiveRequestId,
	iAccept bool,
	iMessage string,
) error {
	pendingRequests, err := c.pendingMaterialReceiveRequestRepositoryService.FetchInboundReceiveMaterialRequestsByUser(
		iContext,
		iUserId,
		[]models.MaterialReceiveRequestStatus{models.PENDING},
	)
	if err != nil {
		return err
	}

	pendingRequest := models.InboundMaterialReceiveRequest{}
	found := false
	for i := range pendingRequests {
		if pendingRequests[i].Id == iRequestId {
			pendingRequest = pendingRequests[i]
			found = true
			break
		}
	}

	if !found {
		return common.NotFound
	}

	acknowledgement, err := c.materialReceiveRequestAcknowledgementRepository.FetchOutboundAcknowledementByRequestId(
		iContext,
		iRequestId,
	)

	if err != nil {
		return err
	}

	if !iAccept {
		/// set request status to rejected
		/// and send response to peer
		err := c.pendingMaterialReceiveRequestRepositoryService.UpdateMaterialReceiveRequestStatus(
			iContext,
			iRequestId,
			models.REJECTED,
		)

		if err != nil {
			return err
		}

		response := peer_material_services.MakeReceiveMaterialResponse(
			string(acknowledgement.ResponseId),
			false,
			iMessage,
			"",
		)

		/// send response to peer
		responseSent := false
		for _, endpoint := range pendingRequest.SenderEndpoints {
			/// try sending
			client, err := c.peerClientFactory.BuildPeerMaterialClientService(
				endpoint.Protocol.Name,
				endpoint.Protocol.MajorVersion,
				endpoint.Protocol.MinorVersion,
				endpoint.Url,
			)

			if err != nil {
				return nil
			}

			_, err = client.SendReceiveMaterialResponse(
				iContext,
				response,
			)

			if err != nil {
				continue
			}

			responseSent = true
			break
		}

		if !responseSent {
			/// temporary code
			c.pendingMaterialReceiveRequestRepositoryService.UpdateMaterialReceiveRequestStatus(
				iContext,
				iRequestId,
				models.PENDING,
			)
			return fmt.Errorf("could not send response")
		}

	}

	return nil
}

func (c *PeerMaterialController) HandleReponse(
	iContext context.Context,
	iRequest peer_material_services.ReceiveMaterialResponse,
) (peer_material_services.ReceiveMaterialResponseAcknowledgement, error) {

	ack, err := c.materialReceiveRequestAcknowledgementRepository.FetchInboundAcknowledementByResponseId(
		iContext,
		models.MaterialReceiveRequestResponseId(iRequest.ResponseId),
	)

	if err != nil {
		return peer_material_services.ReceiveMaterialResponseAcknowledgement{}, err
	}

	err = c.pendingMaterialReceiveRequestRepositoryService.UpdateMaterialReceiveRequestStatus(
		iContext,
		ack.RequestId,
		models.REJECTED,
	)

	if err != nil {
		return peer_material_services.ReceiveMaterialResponseAcknowledgement{}, err
	}

	return peer_material_services.MakeReceiveMaterialResponseAcknowledgement(), nil
}
