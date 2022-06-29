package peer_material_services

import "context"

type PeerMaterialClientServiceI interface {
	SendReceiveMaterialRequest(
		iCtx context.Context,
		iRequest ReceiveMaterialRequestRequest,
	) (ReceiveMaterialRequestResponse, error)
	SendReceiveMaterialResponse(
		iCtx context.Context,
		iResponse ReceiveMaterialResponse,
	) (ReceiveMaterialResponseAcknowledgement, error)
}
