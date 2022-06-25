package models

type SignatureOption struct {
	Id        string `json:"Id"`
	Signature string `json:"Signature"`
}

type MaterialReceiveRequestStatus int

const (
	PENDING  MaterialReceiveRequestStatus = 0
	ACCEPTED                              = 1
	REJECTED                              = 2
)

type MaterialReceiveRequestId int
type MaterialReceiveRequest struct {
	Id                   MaterialReceiveRequestId     `json:"Id"`
	ToBeReceivedMaterial Material                     `json:"ToBeReceivedMaterial"`
	RelatedMaterials     []Material                   `json:"Materials"`
	SignatureOptions     []SignatureOption            `json:"SignatureOption"`
	TransferTime         CustomTime                   `json:"TransferTime"`
	Status               MaterialReceiveRequestStatus `json:"Status"`
}

type OutboundMaterialReceiveRequest struct {
	MaterialReceiveRequest
	RecipientId       PeerId      `json:"RecipientId"`
	SenderUserId      UserId      `json:"SenderUserId"`
	SenderPublicKeyId PublicKeyId `json:"SenderPublicKeyId"`
}

type InboundMaterialReceiveRequest struct {
	MaterialReceiveRequest
	RecipientUserId   UserId      `json:"RecipientUserId"`
	SenderPublicKeyId PublicKeyId `json:"SenderPublicKeyId"`
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

func MakeMaterialReceiveRequest(
	iId MaterialReceiveRequestId,
	iToBeReceivedMaterial Material,
	iRelatedMaterials []Material,
	iSignatureOptions []SignatureOption,
	iTransferTime CustomTime,
	iStatus MaterialReceiveRequestStatus,
) MaterialReceiveRequest {
	return MaterialReceiveRequest{
		Id:                   iId,
		ToBeReceivedMaterial: iToBeReceivedMaterial,
		RelatedMaterials:     iRelatedMaterials,
		SignatureOptions:     iSignatureOptions,
		TransferTime:         iTransferTime,
		Status:               iStatus,
	}
}

func MakeOutboundMaterialReceiveRequest(
	iRequest MaterialReceiveRequest,
	iRecipientId PeerId,
	iSenderUserId UserId,
	iSenderPublicKeyId PublicKeyId,
) OutboundMaterialReceiveRequest {
	return OutboundMaterialReceiveRequest{
		MaterialReceiveRequest: iRequest,
		RecipientId:            iRecipientId,
		SenderUserId:           iSenderUserId,
		SenderPublicKeyId:      iSenderPublicKeyId,
	}
}

func MakeInboundMaterialReceiveRequest(
	iRequest MaterialReceiveRequest,
	iRecipientUserId UserId,
	iSenderPublicKeyId PublicKeyId,
) InboundMaterialReceiveRequest {
	return InboundMaterialReceiveRequest{
		MaterialReceiveRequest: iRequest,
		RecipientUserId:        iRecipientUserId,
		SenderPublicKeyId:      iSenderPublicKeyId,
	}
}
