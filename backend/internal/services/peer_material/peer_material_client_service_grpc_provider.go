package peer_material_services

import "go.uber.org/dig"

func ProvidePeerMaterialClientServiceGrpc(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePeerMaterialClientServiceGrpc, dig.As(new(PeerMaterialClientServiceI)))
}
