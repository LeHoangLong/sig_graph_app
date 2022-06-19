package peer_material_repositories

import (
	"backend/internal/models"
	"context"
)

type SimplifiedPendingReceiveMaterialRequest struct {
	Id                     int
	RecipientId            int
	ToBeReceivedMaterialId int
	RelatedMaterialIds     []int
	SignatureOptions       []models.SignatureOption
	SenderPublicKey        string
	TransferTime           models.CustomTime
}

func makeSimplifiedPendingReceiveMaterialRequest(
	iId int,
	iRecipientId int,
	iToBeReceivedMaterialId int,
	iRelatedMaterialIds []int,
	iSignatureOptions []models.SignatureOption,
	iSenderPublicKey string,
	iTransferTime models.CustomTime,
) SimplifiedPendingReceiveMaterialRequest {
	return SimplifiedPendingReceiveMaterialRequest{
		Id:                     iId,
		RecipientId:            iRecipientId,
		ToBeReceivedMaterialId: iToBeReceivedMaterialId,
		RelatedMaterialIds:     iRelatedMaterialIds,
		SignatureOptions:       iSignatureOptions,
		SenderPublicKey:        iSenderPublicKey,
		TransferTime:           iTransferTime,
	}
}

type PendingMaterialReceiveRequestRepositoryI interface {
	CreatePendingReceiveMaterialRequest(
		iContext context.Context,
		iUser models.User,
		iToBeReceivedMaterial models.Material,
		iRelatedMaterials []models.Material,
		iOptions []models.SignatureOption,
		iSenderPublicKey string,
		iTransferTime models.CustomTime,
		iIsOutbound bool,
	) (models.PendingMaterialReceiveRequest, error)

	FetchPendingReceiveMaterialRequestsByUserId(
		iContext context.Context,
		iUserId int,
		iIsOutbound bool,
	) ([]SimplifiedPendingReceiveMaterialRequest, error)
}
