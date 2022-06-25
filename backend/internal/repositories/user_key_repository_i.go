package repositories

import (
	"backend/internal/models"
	"context"
)

type UserKeyRepositoryI interface {
	FetchUserKeyPairByUser(iUserId models.UserId) ([]models.UserKeyPair, error)
	FetchPublicKeyByPeerId(iContext context.Context, iPeerId models.PeerId) ([]models.PublicKey, error)
	FetchDefaultUserKeyPair(iUserId models.UserId) (models.UserKeyPair, error)
}
