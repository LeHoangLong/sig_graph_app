package graphql

import (
	"backend/internal/models"
	"time"
)

func ParseMaterial(iMaterial models.Material) Material {
	return Material{
		ID:          iMaterial.Id,
		Name:        iMaterial.Name,
		Unit:        iMaterial.Unit,
		Quantity:    iMaterial.Quantity.String(),
		CreatedTime: time.Time(iMaterial.CreatedTime),
	}
}

func ParsePublicKey(iKey models.PublicKey) PublicKey {
	return PublicKey{
		ID:    iKey.Id,
		Value: iKey.Value,
	}
}

func ParsePeer(iPeer models.Peer) Peer {
	keys := make([]*PublicKey, len(iPeer.PublicKey))
	for index, key := range iPeer.PublicKey {
		parsedKey := ParsePublicKey(key)
		keys[index] = &parsedKey
	}
	return Peer{
		ID:         iPeer.Id,
		Alias:      iPeer.Alias,
		PublicKeys: keys,
	}
}
