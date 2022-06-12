package node_contract

import "go.uber.org/dig"

func ProvideSignatureService(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeSignatureService)
}
