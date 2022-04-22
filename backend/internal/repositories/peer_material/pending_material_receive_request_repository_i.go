package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type PendingMaterialReceiveRequestRepositoryI interface {
	createPendingReceiveMaterialRequest(
		iContext context.Context,
		iNodeIds []string,
		iTransferTime models.CustomTime,
		iOptions []models.SignatureOption,
	) (models.PendingMaterialReceiveRequest, error)
}
