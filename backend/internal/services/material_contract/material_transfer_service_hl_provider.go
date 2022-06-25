package material_contract_service

import "go.uber.org/dig"

func ProvideMaterialTransferServiceHl(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(
		MakeMaterialTransferServiceHl,
		dig.As(new(MaterialTransferServiceI)),
	)
}
