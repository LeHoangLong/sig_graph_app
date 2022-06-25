package repositories

import (
	"backend/internal/models"
	"context"
)

type UserRepositoryI interface {
	CreateUser(username string, passwordHash string) (*models.User, error)
	GetUser(username string) (*models.User, error)
	GetUserById(iContext context.Context, iId models.UserId) (models.User, error)
	FindUserWithPublicKey(iPublicKey string) (models.User, error)
	DoesUserExist(username string) (bool, error)
}
