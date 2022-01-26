package repositories

import (
	"backend/internal/common"
	"backend/internal/models"
)

type UserResult struct {
	Error common.WrappedError
	User  *models.User
}

type UsernameResult struct {
	Error    common.WrappedError
	Username string
}

type UserRepositoryI interface {
	/// Does not do any checking if username already exists
	CreateUser(username string, passwordHash string, done chan common.WrappedError)
	/// Returns NotFound if no user with username found
	GetUser(username string, user chan UserResult)
	/// Does not do any checking if username exists
	SetCurrentUsername(username string, done chan common.WrappedError)
	/// Returns NotFound if no current user is set
	GetCurrentUsername(result chan UsernameResult)
}
