package material_contract_service

import (
	"backend/internal/drivers"
	"backend/internal/models"
	graph_id_service "backend/internal/services/graph_id"
	"backend/internal/services/node_contract"
)

type MaterialServiceFactory struct {
	driver      drivers.SmartContractDriverI
	idGenerator graph_id_service.IdGeneratorServiceI
}

func MakeMaterialServiceFactory(
	iDriver drivers.SmartContractDriverI,
	iIdGenerator graph_id_service.IdGeneratorServiceI,
) MaterialServiceFactory {
	return MaterialServiceFactory{
		driver:      iDriver,
		idGenerator: iIdGenerator,
	}
}

func (f MaterialServiceFactory) BuildMaterialCreateService(
	iKey models.UserKeyPair,
) MaterialCreateServiceI {
	signatureService := node_contract.MakeSignatureService()
	return MakeMaterialCreateServiceHl(
		f.driver,
		signatureService,
		f.idGenerator,
	)
}

func (f MaterialServiceFactory) BuildMaterialFetchService() MaterialFetchServiceI {
	return MakeMaterialFetchServiceHl(
		f.driver,
	)
}
