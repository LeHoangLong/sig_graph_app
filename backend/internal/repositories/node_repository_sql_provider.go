package repositories

import "go.uber.org/dig"

func ProvideNodeRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeNodeRepositorySql, dig.As(new(NodeRepositoryI)))
}
