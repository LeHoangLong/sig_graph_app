package repositories

import "go.uber.org/dig"

func ProvideUserEndpointRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(
		MakeUserEndpointRepositorySql,
		dig.As(new(UserEndpointRepositoryI)),
	)
}
