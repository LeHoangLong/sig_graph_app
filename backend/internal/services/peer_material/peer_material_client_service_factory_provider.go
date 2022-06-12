package peer_material_services

import "go.uber.org/dig"

func ProvidePeerMaterialClientService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePeerMaterialClientServiceFactory)
}
