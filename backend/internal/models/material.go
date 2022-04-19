package models

/// Id may be nil if it is obtained from blockchain
type Material struct {
	Node
	Name     string        `json:"Name"`
	Quantity CustomDecimal `json:"Quantity"`
	Unit     string        `json:"Unit"`
}

func NewMaterial(
	iNode Node,
	iName string,
	iQuantity CustomDecimal,
	iUnit string,
) Material {
	return Material{
		Node:     iNode,
		Name:     iName,
		Quantity: iQuantity,
		Unit:     iUnit,
	}
}
