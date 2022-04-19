package peer_material_services

import "context"

type ReceiveMaterialRequestReceivedHandler func(
	iContext context.Context,
	iRequest ReceiveMaterialRequestRequest,
) (ReceiveMaterialRequestResponse, error)

type PeerMaterialServerServiceI interface {
	RegisterReceiveMaterialRequestReceivedHandler(
		iHandler ReceiveMaterialRequestReceivedHandler,
	)
}
