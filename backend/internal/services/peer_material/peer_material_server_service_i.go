package peer_material_services

import "context"

type ReceiveMaterialRequestReceivedHandler interface {
	Handle(
		iContext context.Context,
		iRequest ReceiveMaterialRequestRequest,
	) (ReceiveMaterialRequestResponse, error)
}

type PeerMaterialServerServiceI interface {
}
