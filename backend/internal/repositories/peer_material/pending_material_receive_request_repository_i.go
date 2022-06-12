package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type PendingMaterialReceiveRequestRepositoryI interface {
	CreatePendingReceiveMaterialRequest(
		iContext context.Context,
		iUser models.User,
		iToBeReceivedMaterial models.Material,
		iRelatedMaterials []models.Material,
		iOptions []models.SignatureOption,
		iSenderPublicKey string,
		iTransferTime models.CustomTime,
		iIsOutbound bool,
	) (models.PendingMaterialReceiveRequest, error)
}
