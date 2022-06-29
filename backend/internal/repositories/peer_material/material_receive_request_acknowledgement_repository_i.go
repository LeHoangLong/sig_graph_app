package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type MaterialReceiveRequestAcknowledgementRepositoryI interface {
	SaveAcknowledgement(
		iRequestId models.MaterialReceiveRequestId,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
	/// returns NotFound if no acknowledgement found
	FetchAcknowledementByRequestId(
		iContext context.Context,
		iRequestId models.MaterialReceiveRequestId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
	/// returns NotFound if no acknowledgement found
	FetchAcknowledementByResponseId(
		iContext context.Context,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
}
