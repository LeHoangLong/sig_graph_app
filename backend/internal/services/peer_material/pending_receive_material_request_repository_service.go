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
