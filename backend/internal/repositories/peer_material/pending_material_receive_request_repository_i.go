package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type PendingMaterialReceiveRequestRepositoryI interface {
	createPendingReceiveMaterialRequest(
		iContext context.Context,
		iNodeId string,
		iInfo models.PendingMaterialInfo,
		iOptions []models.SignatureOption,
	) (models.PendingMaterialReceiveRequest, error)
}
