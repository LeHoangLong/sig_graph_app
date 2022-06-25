package graphql

import (
	"backend/internal/models"
	"fmt"
	"time"
)

func ParseOutboundReceiveMaterialRequestRequest(
	iRequest models.OutboundMaterialReceiveRequest,
) ReceiveMaterialRequestRequest {
	parsedMainMaterial := ParseMaterial(iRequest.ToBeReceivedMaterial)
	exposedMaterials := []*Material{}
	for _, material := range iRequest.RelatedMaterials {
		parsedMaterial := ParseMaterial(material)
		exposedMaterials = append(exposedMaterials, &parsedMaterial)
	}

	return ReceiveMaterialRequestRequest{
		TransferMaterial:  &parsedMainMaterial,
		ExposedMaterials:  exposedMaterials,
		TransferTime:      time.Time(iRequest.TransferTime),
		SenderPublicKeyID: int(iRequest.SenderPublicKeyId),
		ID:                int(iRequest.Id),
		Status:            ParseStatus(iRequest.Status),
	}
}

func ParseInboundReceiveMaterialRequestRequest(
	iRequest models.InboundMaterialReceiveRequest,
) ReceiveMaterialRequestRequest {
	parsedMainMaterial := ParseMaterial(iRequest.ToBeReceivedMaterial)
	exposedMaterials := []*Material{}
	for _, material := range iRequest.RelatedMaterials {
		parsedMaterial := ParseMaterial(material)
		exposedMaterials = append(exposedMaterials, &parsedMaterial)
	}

	return ReceiveMaterialRequestRequest{
		TransferMaterial:  &parsedMainMaterial,
		ExposedMaterials:  exposedMaterials,
		TransferTime:      time.Time(iRequest.TransferTime),
		SenderPublicKeyID: int(iRequest.SenderPublicKeyId),
		ID:                int(iRequest.Id),
		Status:            ParseStatus(iRequest.Status),
	}
}

func ParseMaterial(iMaterial models.Material) Material {
	previousNodeHashedIds := []*string{}
	for previousNodeHashedId := range iMaterial.PreviousNodeHashedIds {
		previousNodeHashedIds = append(previousNodeHashedIds, &previousNodeHashedId)
	}
	nextNodeHashedIds := []*string{}
	for nextNodeHashedId := range iMaterial.PreviousNodeHashedIds {
		nextNodeHashedIds = append(nextNodeHashedIds, &nextNodeHashedId)
	}
	return Material{
		ID:                     int(*iMaterial.Id),
		NodeID:                 iMaterial.NodeId,
		Name:                   iMaterial.Name,
		Unit:                   iMaterial.Unit,
		Quantity:               iMaterial.Quantity.String(),
		CreatedTime:            time.Time(iMaterial.CreatedTime),
		OwnerPublicKey:         iMaterial.OwnerPublicKey.Value,
		PreviousNodesHashedIds: previousNodeHashedIds,
		NextNodesHashedIds:     nextNodeHashedIds,
	}
}

func ParsePublicKey(iKey models.PublicKey) (PublicKey, error) {
	if iKey.Id != nil {
		return PublicKey{
			ID:    int(*iKey.Id),
			Value: iKey.Value,
		}, nil
	} else {
		return PublicKey{}, fmt.Errorf("public key not yet saved to database")
	}
}

func ParsePeer(iPeer models.Peer) (Peer, error) {
	keys := make([]*PublicKey, len(iPeer.PublicKey))
	for index, key := range iPeer.PublicKey {
		parsedKey, err := ParsePublicKey(key)
		if err != nil {
			return Peer{}, err
		}
		keys[index] = &parsedKey
	}
	return Peer{
		ID:         int(iPeer.Id),
		Alias:      iPeer.Alias,
		PublicKeys: keys,
	}, nil
}

func CompileStatus(iStatus ReceiveMaterialRequestRequestStatus) models.MaterialReceiveRequestStatus {
	switch iStatus {
	case ReceiveMaterialRequestRequestStatusPending:
		return models.PENDING
	case ReceiveMaterialRequestRequestStatusAccepted:
		return models.ACCEPTED
	case ReceiveMaterialRequestRequestStatusRejected:
		return models.REJECTED
	}
	panic("Unsupported status")
}

func ParseStatus(iStatus models.MaterialReceiveRequestStatus) ReceiveMaterialRequestRequestStatus {
	switch iStatus {
	case models.PENDING:
		return ReceiveMaterialRequestRequestStatusPending
	case models.ACCEPTED:
		return ReceiveMaterialRequestRequestStatusAccepted
	case models.REJECTED:
		return ReceiveMaterialRequestRequestStatusRejected
	}
	panic("Unsupported status")
}
