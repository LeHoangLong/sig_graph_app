package repositories

import (
	"backend/internal/models"
	"context"
)

type PeerKeyRepositoryI interface {
	CreateOrFetchPeerKeysByValue(
		iContext context.Context,
		iOwnerId models.UserId,
		iPeerKeys []string,
	) ([]models.PeerKey, error)
	FetchPublicKeysById(
		iContext context.Context,
		iKeysId map[models.PublicKeyId]bool,
	) (map[models.PublicKeyId]models.PublicKey, error)
}
