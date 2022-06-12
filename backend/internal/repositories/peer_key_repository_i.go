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
}
