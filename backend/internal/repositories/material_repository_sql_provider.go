package repositories

import "go.uber.org/dig"

func ProvideMaterialRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeMaterialRepositorySql, dig.As(new(MaterialRepositoryI)))
}
