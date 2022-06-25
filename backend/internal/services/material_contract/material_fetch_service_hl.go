package material_contract_service

import (
	"backend/internal/drivers"
	"backend/internal/models"
	"encoding/json"
)

type MaterialFetchServiceHl struct {
	driver drivers.SmartContractDriverI
}

func MakeMaterialFetchServiceHl(
	iDriver drivers.SmartContractDriverI,
) MaterialFetchServiceHl {
	return MaterialFetchServiceHl{
		iDriver,
	}
}

func (s MaterialFetchServiceHl) GetMaterialById(
	iNodeId string,
) (models.Material, error) {
	nodeJson, err := s.driver.Query(
		"GetMaterial",
		iNodeId,
	)
	if err != nil {
		return models.Material{}, err
	}

	var materialSc material
	err = json.Unmarshal(nodeJson, &materialSc)
	if err != nil {
		return models.Material{}, err
	}

	node := models.MakeNode(
		nil,
		materialSc.Id,
		materialSc.IsFinalized,
		materialSc.PreviousNodeHashedIds,
		materialSc.NextNodeHashedIds,
		map[models.NodeId]bool{},
		map[models.NodeId]bool{},
		models.MakePublicKey(
			nil,
			materialSc.OwnerPublicKey,
		),
		materialSc.CreatedTime,
		materialSc.Signature,
		"material",
	)

	return models.NewMaterial(
		node,
		materialSc.Name,
		materialSc.Quantity,
		materialSc.Unit,
	), nil
}
