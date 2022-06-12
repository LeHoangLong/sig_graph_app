package peer_material_repositories

import "backend/internal/models"

type PendingMaterialReceiveRequestResponseRepositoryI interface {
	SaveResponse(
		iRequestId int,
		iResponseId string,
	) (models.PendingMaterialReceiveRequestResponse, error)
}
