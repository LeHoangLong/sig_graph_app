package models

type SignatureOption struct {
	Id        string `json:"Id"`
	Signature string `json:"Signature"`
}

type PendingMaterialReceiveRequest struct {
	Id                   int               `json:"Id"`
	RecipientId          int               `json:"RecipientId"`
	ToBeReceivedMaterial Material          `json:"ToBeReceivedMaterial"`
	RelatedMaterials     []Material        `json:"Materials"`
	SignatureOptions     []SignatureOption `json:"SignatureOption"`
	SenderPublicKey      string            `json:"SenderPublicKey"`
	TransferTime         CustomTime        `json:"TransferTime"`
}

func MakeSignatureOption(
	iId string,
	iSignature string,
) SignatureOption {
	return SignatureOption{
		Id:        iId,
		Signature: iSignature,
	}
}

func MakePendingMaterialReceiveRequest(
	iId int,
	iRecipientId int,
	iToBeReceivedMaterial Material,
	iRelatedMaterials []Material,
	iSignatureOptions []SignatureOption,
	iSenderPublicKey string,
	iTransferTime CustomTime,
) PendingMaterialReceiveRequest {
	return PendingMaterialReceiveRequest{
		Id:                   iId,
		RecipientId:          iRecipientId,
		ToBeReceivedMaterial: iToBeReceivedMaterial,
		RelatedMaterials:     iRelatedMaterials,
		SignatureOptions:     iSignatureOptions,
		SenderPublicKey:      iSenderPublicKey,
		TransferTime:         iTransferTime,
	}
}
