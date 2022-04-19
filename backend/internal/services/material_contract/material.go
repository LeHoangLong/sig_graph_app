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
	return material{
		Name:     iMaterial.Name,
		Unit:     iMaterial.Unit,
		Quantity: iMaterial.Quantity,
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
