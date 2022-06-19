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
	idHasher                                        node_contract.IdHasherI
	pendingMaterialReceiveRequestRepositoryService  *peer_material_services.PendingReceiveMaterialRequestRepositoryService
	peerClientFactory                               peer_material_services.PeerMaterialClientServiceFactory
	userKeyRepository                               repositories.UserKeyRepositoryI
	optionGenerator                                 material_contract_service.SignatureOptionGenerator
	materialRepositoryService                       services.MaterialRepositoryService
	pendingMaterialReceiveRequestResponseRepository peer_material_repositories.PendingMaterialReceiveRequestResponseRepositoryI
	peerRepository                                  repositories.PeerRepositoryI
	peerKeyRepository                               repositories.PeerKeyRepositoryI
}

func MakePeerMaterialController(
	iUserService services.UserService,
	iMaterialFetchService material_contract_service.MaterialFetchServiceI,
	iIdHasher node_contract.IdHasherI,
	iPendingMaterialReceiveRequestRepositoryService *peer_material_services.PendingReceiveMaterialRequestRepositoryService,
	iPeerClientFactory peer_material_services.PeerMaterialClientServiceFactory,
	iUserKeyRepository repositories.UserKeyRepositoryI,
	iOptionGenerator material_contract_service.SignatureOptionGenerator,
	iMaterialRepositoryService services.MaterialRepositoryService,
	iPendingMaterialReceiveRequestResponseRepository peer_material_repositories.PendingMaterialReceiveRequestResponseRepositoryI,
	iPeerRepository repositories.PeerRepositoryI,
	iPeerKeyRepository repositories.PeerKeyRepositoryI,
) *PeerMaterialController {
	return &PeerMaterialController{
		userService:          iUserService,
		materialFetchService: iMaterialFetchService,
		idHasher:             iIdHasher,
		pendingMaterialReceiveRequestRepositoryService: iPendingMaterialReceiveRequestRepositoryService,
		peerClientFactory:         iPeerClientFactory,
		userKeyRepository:         iUserKeyRepository,
		optionGenerator:           iOptionGenerator,
		materialRepositoryService: iMaterialRepositoryService,
		pendingMaterialReceiveRequestResponseRepository: iPendingMaterialReceiveRequestResponseRepository,
		peerRepository:    iPeerRepository,
		peerKeyRepository: iPeerKeyRepository,
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
		user,
		materialOwnerKeys,
	)
	materialOwnerKeyModelsIdMap := map[string]int{}
	for index := range materialOwnerKeyModels {
		materialOwnerKeyModelsIdMap[materialOwnerKeyModels[index].Value] = *materialOwnerKeyModels[index].Id
	}

	for index := range fetchedMaterials {
		id := materialOwnerKeyModelsIdMap[fetchedMaterials[index].OwnerPublicKey.Value]
		fetchedMaterials[index].OwnerPublicKey.Id = &id
	}

	fetchedMaterials, err = c.materialRepositoryService.SaveMaterials(
		iContext,
		fetchedMaterials,
	)
	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, err
	}

	options := []models.SignatureOption{}
	for _, option := range iRequest.SignatureOptions {
		options = append(options, models.MakeSignatureOption(option.NodeId, option.Signature))
	}
	pendingRequest, err := c.pendingMaterialReceiveRequestRepositoryService.SavePendingReceiveMaterialRequest(
		iContext,
		user,
		fetchedMaterials[0],
		fetchedMaterials[1:],
		options,
		iRequest.SenderPublicKey,
		models.CustomTime(iRequest.TransferTime),
		false,
	)

	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, nil
	}

	response, err := c.pendingMaterialReceiveRequestResponseRepository.SaveResponse(
		pendingRequest.Id,
		strconv.Itoa(pendingRequest.Id),
	)

	if err != nil {
		return peer_material_services.ReceiveMaterialRequestResponse{}, nil
	}

	return peer_material_services.MakeReceiveMaterialRequestResponse(
		response.ResponseId,
		true,
	), nil
}

func (c *PeerMaterialController) SendRequest(
	iContext context.Context,
	iPeerId int,
	iRecipientPublicKeyId int,
	iMainMaterialId int,
	iRelatedMaterialId []int,
) (models.PendingMaterialReceiveRequest, error) {
	senderId, err := services.GetCurrentUserFromContext(iContext)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	sender, err := c.userService.GetUserById(iContext, senderId)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	keys, err := c.userKeyRepository.FetchPublicKeyByPeerId(iContext, iPeerId)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	recipientPublicKey := ""
	for _, key := range keys {
		if *key.Id == iRecipientPublicKeyId {
			recipientPublicKey = key.Value
		}
	}

	if recipientPublicKey == "" {
		return models.PendingMaterialReceiveRequest{}, common.NotFound
	}

	transferTime := models.CustomTime(time.Now())
	key, err := c.userKeyRepository.FetchDefaultUserKeyPair(sender.Id)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	mainMaterial, relatedMaterials, err := c.materialRepositoryService.FetchMaterialsAndRelated(iContext, iMainMaterialId)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
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
		return models.PendingMaterialReceiveRequest{}, err
	}

	request, err := c.pendingMaterialReceiveRequestRepositoryService.SavePendingReceiveMaterialRequest(
		iContext,
		sender,
		mainMaterial,
		filteredRelatedMaterials,
		options,
		key.PublicKey.Value,
		transferTime,
		true,
	)

	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
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

	peerRequest := peer_material_services.MakeReceiveMaterialRequestRequest(
		recipientPublicKey,
		mainMaterial.NodeId,
		nodesMap,
		time.Time(transferTime),
		request.SenderPublicKey,
		serviceLayerOptions,
	)

	endpoints, err := c.peerRepository.FetchPeerEndPoints(iContext, iPeerId)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	requestReceived := false
	for _, endpoint := range endpoints {
		client, err := c.peerClientFactory.BuildPeerMaterialClientService(endpoint)
		if err != nil {
			continue
		}

		response, err := client.SendReceiveMaterialRequest(iContext, peerRequest)
		if err != nil {
			return models.PendingMaterialReceiveRequest{}, err
		}

		if !response.IsRequestAcknowledged {
			return models.PendingMaterialReceiveRequest{}, fmt.Errorf("Request rejected")
		}

		_, err = c.pendingMaterialReceiveRequestResponseRepository.SaveResponse(
			request.Id,
			response.ResponseId,
		)

		requestReceived = true
		break
	}

	if !requestReceived {
		return models.PendingMaterialReceiveRequest{}, errors.New("could not send request")
	}

	return request, nil
}

func (c *PeerMaterialController) FetchReceivedPendingMaterialReceiveRequests(
	iContext context.Context,
	iSenderId int,
) ([]models.PendingMaterialReceiveRequest, error) {
	return c.pendingMaterialReceiveRequestRepositoryService.FetchPendingReceiveMaterialRequestsByUser(
		iContext,
		iSenderId,
		false,
	)
}
