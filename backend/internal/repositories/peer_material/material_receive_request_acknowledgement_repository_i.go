package peer_material_repositories

import "backend/internal/models"

type MaterialReceiveRequestAcknowledgementRepositoryI interface {
	SaveAcknowledgement(
		iRequestId models.MaterialReceiveRequestId,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
}
