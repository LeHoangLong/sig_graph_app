package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type MaterialReceiveRequestAcknowledgementRepositoryI interface {
	SaveOutboundAcknowledgement(
		iRequestId models.MaterialReceiveRequestId,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
	/// returns NotFound if no acknowledgement found
	FetchOutboundAcknowledementByRequestId(
		iContext context.Context,
		iRequestId models.MaterialReceiveRequestId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
	/// returns NotFound if no acknowledgement found
	FetchOutboundAcknowledementByResponseId(
		iContext context.Context,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)

	SaveInboundAcknowledgement(
		iRequestId models.MaterialReceiveRequestId,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
	/// returns NotFound if no acknowledgement found
	FetchInboundAcknowledementByRequestId(
		iContext context.Context,
		iRequestId models.MaterialReceiveRequestId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
	/// returns NotFound if no acknowledgement found
	FetchInboundAcknowledementByResponseId(
		iContext context.Context,
		iResponseId models.MaterialReceiveRequestResponseId,
	) (models.MaterialReceiveRequestAcknowledgement, error)
}
