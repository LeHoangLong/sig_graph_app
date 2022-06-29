package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type SimplifiedInboundReceiveMaterialRequest struct {
	Id                     models.MaterialReceiveRequestId
	RecipientUserId        models.UserId
	ToBeReceivedMaterialId models.NodeId
	RelatedMaterialIds     []models.NodeId
	SignatureOptions       []models.SignatureOption
	SenderPublicKey        models.PublicKeyId
	TransferTime           models.CustomTime
	Status                 models.MaterialReceiveRequestStatus
	SenderEndpoints        []models.SenderEndpoint
}

func makeSimplifiedInboundReceiveMaterialRequest(
	iId models.MaterialReceiveRequestId,
	iRecipientUserId models.UserId,
	iToBeReceivedMaterialId models.NodeId,
	iRelatedMaterialIds []models.NodeId,
	iSignatureOptions []models.SignatureOption,
	iSenderPublicKey models.PublicKeyId,
	iTransferTime models.CustomTime,
	iStatus models.MaterialReceiveRequestStatus,
	iSenderEndpoints []models.SenderEndpoint,
) SimplifiedInboundReceiveMaterialRequest {
	return SimplifiedInboundReceiveMaterialRequest{
		Id:                     iId,
		RecipientUserId:        iRecipientUserId,
		ToBeReceivedMaterialId: iToBeReceivedMaterialId,
		RelatedMaterialIds:     iRelatedMaterialIds,
		SignatureOptions:       iSignatureOptions,
		SenderPublicKey:        iSenderPublicKey,
		TransferTime:           iTransferTime,
		Status:                 iStatus,
		SenderEndpoints:        iSenderEndpoints,
	}
}

type MaterialReceiveRequestRepositoryI interface {
	CreateOutboundReceiveMaterialRequest(
		iContext context.Context,
		iRecipientPeerId models.PeerId,
		iSenderUserId models.UserId,
		iSenderPublicKeyId models.PublicKeyId,
		iToBeReceivedMaterial models.Material,
		iRelatedMaterials []models.Material,
		iOptions []models.SignatureOption,
		iTransferTime models.CustomTime,
		iStatus models.MaterialReceiveRequestStatus,
	) (models.OutboundMaterialReceiveRequest, error)

	CreateInboundReceiveMaterialRequest(
		iContext context.Context,
		iRecipientUserId models.UserId,
		iSenderPublicKeyId models.PublicKeyId,
		iToBeReceivedMaterial models.Material,
		iRelatedMaterials []models.Material,
		iOptions []models.SignatureOption,
		iTransferTime models.CustomTime,
		iStatus models.MaterialReceiveRequestStatus,
		iSenderEndpoints []models.SenderEndpoint,
	) (models.InboundMaterialReceiveRequest, error)

	FetchInboundReceiveMaterialRequestsByUserId(
		iContext context.Context,
		iUserId models.UserId,
		iStatus []models.MaterialReceiveRequestStatus,
	) ([]SimplifiedInboundReceiveMaterialRequest, error)

	/// returns NotFound if iRequestId not found
	UpdateMaterialReceiveRequestStatus(
		iContext context.Context,
		iRequestId models.MaterialReceiveRequestId,
		iStatus models.MaterialReceiveRequestStatus,
	) error
}
