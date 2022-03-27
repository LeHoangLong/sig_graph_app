package repositories

import (
	"backend/internal/models"
)

type UserKeyRepositoryI interface {
	FetchUserKeyPairByUser(iUserId int) ([]models.UserKeyPair, error)
}
