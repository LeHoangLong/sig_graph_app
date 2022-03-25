package repositories

import "backend/internal/models"

type UserRepositoryI interface {
	CreateUser(username string, passwordHash string) (*models.User, error)
	GetUser(username string) (*models.User, error)
	DoesUserExist(username string) error
}
