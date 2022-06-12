package repositories

import "go.uber.org/dig"

func ProvideUserRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(
		MakeUserRepositorySql,
		dig.As(new(UserRepositoryI)),
	)
}
