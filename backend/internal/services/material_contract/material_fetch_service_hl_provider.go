package material_contract_service

import "go.uber.org/dig"

func ProvideMaterialFetchServiceHl(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeMaterialFetchServiceHl, dig.As(new(MaterialFetchServiceI)))
}
