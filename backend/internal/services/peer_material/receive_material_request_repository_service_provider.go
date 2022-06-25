package peer_material_services

import "go.uber.org/dig"

func ProvideReceiveMaterialRequestRepositoryService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeReceiveMaterialRequestRepositoryService)
}
