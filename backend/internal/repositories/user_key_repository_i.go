package repositories

import (
	"backend/internal/models"
	"context"
)

type UserKeyRepositoryI interface {
	FetchUserKeyPairByUser(iUserId int) ([]models.UserKeyPair, error)
	FetchPublicKeyByPeerId(iContext context.Context, iPeerId int) ([]models.PublicKey, error)
	FetchDefaultUserKeyPair(iUserId int) (models.UserKeyPair, error)
}
