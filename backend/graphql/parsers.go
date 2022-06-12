package graphql

import (
	"backend/internal/models"
	"fmt"
	"time"
)

func ParseReceiveMaterialRequestRequest(
	iRequest models.PendingMaterialReceiveRequest,
) ReceiveMaterialRequestRequest {
	parsedMainMaterial := ParseMaterial(iRequest.ToBeReceivedMaterial)
	exposedMaterials := []*Material{}
	for _, material := range iRequest.RelatedMaterials {
		parsedMaterial := ParseMaterial(material)
		exposedMaterials = append(exposedMaterials, &parsedMaterial)
	}

	return ReceiveMaterialRequestRequest{
		TransferMaterial: &parsedMainMaterial,
		ExposedMaterials: exposedMaterials,
		TransferTime:     time.Time(iRequest.TransferTime),
	}
}

func ParseMaterial(iMaterial models.Material) Material {
	return Material{
		ID:          *iMaterial.Id,
		NodeID:      iMaterial.NodeId,
		Name:        iMaterial.Name,
		Unit:        iMaterial.Unit,
		Quantity:    iMaterial.Quantity.String(),
		CreatedTime: time.Time(iMaterial.CreatedTime),
	}
}

func ParsePublicKey(iKey models.PublicKey) (PublicKey, error) {
	if iKey.Id != nil {
		return PublicKey{
			ID:    *iKey.Id,
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
		ID:         iPeer.Id,
		Alias:      iPeer.Alias,
		PublicKeys: keys,
	}, nil
}
