package models

type PendingMaterialInfo struct {
	Name        string        `json:"Name"`
	Quantity    CustomDecimal `json:"Quantity"`
	Unit        string        `json:"Unit"`
	CreatedTime CustomTime    `json:"CreatedTime"`
}

type SignatureOption struct {
	Id        string `json:"Id"`
	Signature string `json:"Signature"`
}

type PendingMaterialReceiveRequest struct {
	Id                  string              `json:"Id"`
	PendingMaterialInfo PendingMaterialInfo `json:"PendingMaterialInfo"`
	NodeId              string              `json:"PreviousNodeId"`
	SignatureOptions    []SignatureOption   `json:"SignatureOption"`
}
