package peer_material_services

import "context"

type ReceiveMaterialRequestReceivedHandler interface {
	Handle(
		iContext context.Context,
		iRequest ReceiveMaterialRequestRequest,
	) (ReceiveMaterialRequestResponse, error)
}

type ReceiveMaterialResponseHandler interface {
	HandleReponse(
		iContext context.Context,
		iRequest ReceiveMaterialResponse,
	) (ReceiveMaterialResponseAcknowledgement, error)
}

type PeerMaterialServerServiceI interface {
}
