package peer_material_services

import "context"

type PeerMaterialClientServiceI interface {
	SendReceiveMaterialRequest(
		iCtx context.Context,
		iRequest ReceiveMaterialRequestRequest,
	) (ReceiveMaterialRequestResponse, error)
}
