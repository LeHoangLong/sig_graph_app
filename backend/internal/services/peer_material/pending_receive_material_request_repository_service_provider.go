package peer_material_services

import "go.uber.org/dig"

func ProvidePendingReceiveMaterialRequestRepositoryService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePendingReceiveMaterialRequestRepositoryService)
}
