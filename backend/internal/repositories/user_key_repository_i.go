package repositories

import (
	"backend/internal/models"
)

type UserKeyRepositoryI interface {
	FetchUserKeyPairByUser(iUserId int) ([]models.UserKeyPair, error)
	FetchDefaultUserKeyPair(iUserId int) (models.UserKeyPair, error)
}
