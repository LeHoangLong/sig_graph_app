package node_contract

import "go.uber.org/dig"

func ProvideIdHasher(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeIdHasher, dig.As(new(IdHasherI)))
}
