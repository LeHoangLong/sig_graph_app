package material_contract_service

import (
	"backend/internal/drivers"
	"backend/internal/models"
	graph_id_service "backend/internal/services/graph_id"
	"backend/internal/services/node_contract"
	"context"
	"time"
)

type MaterialCreateService struct {
	driver           drivers.SmartContractDriverI
	signatureService node_contract.SignatureService
	idGenerator      graph_id_service.IdGeneratorServiceI
}

func MakeMaterialCreateServiceHl(
	iDriver drivers.SmartContractDriverI,
	iSignatureService node_contract.SignatureService,
	iIdGenerator graph_id_service.IdGeneratorServiceI,
) MaterialCreateService {
	return MaterialCreateService{
		driver:           iDriver,
		signatureService: iSignatureService,
		idGenerator:      iIdGenerator,
	}
}

func (c MaterialCreateService) CreateMaterial(
	iContext context.Context,
	iName string,
	iUnit string,
	iQuantity models.CustomDecimal,
	iOwnerKey models.UserKeyPair,
) (models.Material, error) {
	createdTime := models.CustomTime(time.Now())
	nodeId := c.idGenerator.NewId(iContext)
	nodeHeader := node_contract.MakeNodeHeader(
		nodeId,
		false,
		map[string]bool{},
		map[string]bool{},
		iOwnerKey.PublicKey.Value,
		createdTime,
		"",
	)
	materialSc := makeMaterial(
		nodeHeader,
		iName,
		iUnit,
		iQuantity,
	)

	signature, err := c.signatureService.CreateNodeSignature(
		iOwnerKey.PrivateKey,
		&materialSc,
	)

	if err != nil {
		return models.Material{}, err
	}

	_, err = c.driver.CreateTransaction(
		"CreateMaterial",
		nodeId,
		iName,
		iUnit,
		iQuantity.String(),
		iOwnerKey.PublicKey.Value,
		createdTime.String(),
		signature,
	)

	if err != nil {
		return models.Material{}, err
	}

	node := models.MakeNode(
		nil,
		nodeId,
		false,
		map[string]bool{},
		map[string]bool{},
		map[int]bool{},
		map[int]bool{},
		iOwnerKey.PublicKey,
		createdTime,
		signature,
		"material",
	)

	newMaterial := models.NewMaterial(
		node,
		iName,
		iQuantity,
		iUnit,
	)
	return newMaterial, err
}
