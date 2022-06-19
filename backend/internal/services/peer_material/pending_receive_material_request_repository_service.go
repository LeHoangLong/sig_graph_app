package peer_material_services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	peer_material_repositories "backend/internal/repositories/peer_material"
	"context"
)

type PendingReceiveMaterialRequestRepositoryService struct {
	materialRepository                      repositories.MaterialRepositoryI
	pendingReceiveMaterialRequestRepository peer_material_repositories.PendingMaterialReceiveRequestRepositoryI
}

func MakePendingReceiveMaterialRequestRepositoryService(
	iMaterialRepository repositories.MaterialRepositoryI,
	iPendingReceiveMaterialRequestRepository peer_material_repositories.PendingMaterialReceiveRequestRepositoryI,
) *PendingReceiveMaterialRequestRepositoryService {
	return &PendingReceiveMaterialRequestRepositoryService{
		materialRepository:                      iMaterialRepository,
		pendingReceiveMaterialRequestRepository: iPendingReceiveMaterialRequestRepository,
	}
}

func (s *PendingReceiveMaterialRequestRepositoryService) SavePendingReceiveMaterialRequest(
	iContext context.Context,
	iUser models.User,
	iMainMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iSenderPublicKey string,
	iTransferTime models.CustomTime,
	iIsOutbound bool,
) (models.PendingMaterialReceiveRequest, error) {
	return s.pendingReceiveMaterialRequestRepository.CreatePendingReceiveMaterialRequest(
		iContext,
		iUser,
		iMainMaterial,
		iRelatedMaterials,
		iOptions,
		iSenderPublicKey,
		iTransferTime,
		iIsOutbound,
	)
}

func (s *PendingReceiveMaterialRequestRepositoryService) FetchPendingReceiveMaterialRequestsByUser(
	iContext context.Context,
	iUserId int,
	iIsOutbound bool,
) ([]models.PendingMaterialReceiveRequest, error) {
	simplifiedRequests, err := s.pendingReceiveMaterialRequestRepository.FetchPendingReceiveMaterialRequestsByUserId(
		iContext,
		iUserId,
		iIsOutbound,
	)

	if err != nil {
		return []models.PendingMaterialReceiveRequest{}, err
	}

	materialsIdToFetch := map[int]bool{}
	for _, request := range simplifiedRequests {
		materialsIdToFetch[request.ToBeReceivedMaterialId] = true
		for relatedMaterialId := range request.RelatedMaterialIds {
			materialsIdToFetch[relatedMaterialId] = true
		}
	}

	materials, err := s.materialRepository.FetchMaterialsById(
		iContext,
		materialsIdToFetch,
	)

	if err != nil {
		return []models.PendingMaterialReceiveRequest{}, err
	}

	requests := []models.PendingMaterialReceiveRequest{}
	for _, request := range simplifiedRequests {
		relatedMaterials := []models.Material{}
		for _, id := range request.RelatedMaterialIds {
			relatedMaterials = append(relatedMaterials, materials[id])
		}

		pendingRequest := models.MakePendingMaterialReceiveRequest(
			request.Id,
			request.RecipientId,
			materials[request.ToBeReceivedMaterialId],
			relatedMaterials,
			request.SignatureOptions,
			request.SenderPublicKey,
			request.TransferTime,
		)

		requests = append(requests, pendingRequest)
	}

	return requests, nil
}
