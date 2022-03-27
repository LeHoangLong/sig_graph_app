package models

type Material struct {
	Id             string        `json:"ID"`
	Name           string        `json:"Name"`
	Quantity       CustomDecimal `json:"Quantity"`
	Unit           string        `json:"Unit"`
	CreatedTime    CustomTime    `json:"CreatedTime"`
	OwnerPublicKey string        `json:"OwnerPublicKey"`
}

func NewMaterial(
	iId string,
	iName string,
	iQuantity CustomDecimal,
	iUnit string,
	iCreatedTime CustomTime,
	iOwnerPublicKey string,
) Material {
	return Material{
		iId,
		iName,
		iQuantity,
		iUnit,
		iCreatedTime,
		iOwnerPublicKey,
	}
}
