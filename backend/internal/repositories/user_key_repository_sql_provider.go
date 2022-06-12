package repositories

import "go.uber.org/dig"

func ProvideUserKeyRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeUserKeyRepositorySql, dig.As(new(UserKeyRepositoryI)))
}
