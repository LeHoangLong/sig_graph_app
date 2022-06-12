package repositories

import "go.uber.org/dig"

func ProvideMaterialRepositoryFactory(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeMaterialRepositoryFactory)
}
