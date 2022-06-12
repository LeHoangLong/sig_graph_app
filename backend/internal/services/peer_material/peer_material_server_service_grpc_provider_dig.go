package peer_material_services

import "go.uber.org/dig"

func ProvidePeerMaterialServerService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePeerMaterialServerServiceGrpc)
}
