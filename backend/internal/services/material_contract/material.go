package material_contract_service

import (
	"backend/internal/models"
	"backend/internal/services/node_contract"
)

type material struct {
	node_contract.NodeHeader
	Name     string               `json:"Name"`
	Unit     string               `json:"Unit"`
	Quantity models.CustomDecimal `json:"Quantity"`
}

func (m *material) GetHeader() node_contract.NodeHeader {
	return m.NodeHeader
}
func (m *material) SetHeader(iHeader node_contract.NodeHeader) {
	m.NodeHeader = iHeader
}

func makeMaterialFromModel(iMaterial models.Material) material {
	header := node_contract.MakeNodeHeader(
		iMaterial.NodeId,
		iMaterial.IsFinalized,
		iMaterial.PreviousNodeHashedIds,
		iMaterial.NextNodeHashedIds,
		iMaterial.OwnerPublicKey.Value,
		iMaterial.CreatedTime,
		iMaterial.Signature,
	)
	return material{
		NodeHeader: header,
		Name:       iMaterial.Name,
		Unit:       iMaterial.Unit,
		Quantity:   iMaterial.Quantity,
	}
}

func makeMaterial(
	iNodeHeader node_contract.NodeHeader,
	iName string,
	iUnit string,
	iQuantity models.CustomDecimal,
) material {
	return material{
		NodeHeader: iNodeHeader,
		Name:       iName,
		Unit:       iUnit,
		Quantity:   iQuantity,
	}
}
