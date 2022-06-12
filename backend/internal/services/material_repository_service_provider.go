package services

import "go.uber.org/dig"

func ProvideMaterialRepositoryService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeMaterialRepositoryService)
}
