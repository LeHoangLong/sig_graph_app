package material_contract_service

import "go.uber.org/dig"

func ProvideSignatureOptionGenerator(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeSignatureOptionGenerator)
}
