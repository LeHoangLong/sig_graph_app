package models

type Material struct {
	Id          string        `json:"ID"`
	Name        string        `json:"Name"`
	Quantity    CustomDecimal `json:"Quantity"`
	Unit        string        `json:"Unit"`
	CreatedTime CustomTime    `json:"CreatedTime"`
}

func NewMaterial(
	iId string,
	iName string,
	iQuantity CustomDecimal,
	iUnit string,
	iCreatedTime CustomTime,
) Material {
	return Material{
		iId,
		iName,
		iQuantity,
		iUnit,
		iCreatedTime,
	}
}
