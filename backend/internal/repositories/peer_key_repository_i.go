package repositories

import (
	"backend/internal/models"
	"context"
)

type PeerKeyRepositoryI interface {
	CreateOrFetchPeerKeysByValue(
		iContext context.Context,
		iOwner models.User,
		iPeerKeys []string,
	) ([]models.PeerKey, error)
	FetchPublicKeysById(
		iContext context.Context,
		iKeysId map[models.PublicKeyId]bool,
	) (map[models.PublicKeyId]models.PublicKey, error)
}
