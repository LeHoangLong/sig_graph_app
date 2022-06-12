package services

import "go.uber.org/dig"

func ProvideUserService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeUserService)
}
