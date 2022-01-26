package services

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserStorageService struct {
	repository repositories.UserRepositoryI
}

func MakeUserStorageService(
	repository repositories.UserRepositoryI,
) *UserStorageService {
	return &UserStorageService{
		repository: repository,
	}
}

type UserResult struct {
	User  *models.User
	Error common.WrappedError
}

/// Returns NotFound if no username and password match
/// Else returns nil Error
func (s *UserStorageService) VerifyUser(
	username string,
	password string,
	result chan common.WrappedError,
) {
	userResultChan := make(chan repositories.UserResult)
	go s.repository.GetUser(username, userResultChan)
	userResult := <-userResultChan
	if userResult.Error.Error != nil {
		result <- userResult.Error
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userResult.User.PasswordHash), []byte(password))
	if err != nil {
		result <- common.WrappedError{
			Code:  common.NotFound,
			Error: err,
		}
		return
	}

	result <- common.WrappedError{
		Error: nil,
	}
}

func (s *UserStorageService) SetCurrentUser(
	username string,
	result chan common.WrappedError,
) {
	go s.repository.SetCurrentUsername(username, result)
}

func (s *UserStorageService) GetCurrentUser(oResult chan UserResult) {
	currentUsernameResultChan := make(chan repositories.UsernameResult)
	go s.repository.GetCurrentUsername(currentUsernameResultChan)
	currentUsernameResult := <-currentUsernameResultChan
	if currentUsernameResult.Error.Error != nil {
		oResult <- UserResult{
			Error: currentUsernameResult.Error,
			User:  nil,
		}
		return
	}

	userResultChan := make(chan repositories.UserResult)
	go s.repository.GetUser(currentUsernameResult.Username, userResultChan)

	result := <-userResultChan
	oResult <- UserResult{
		Error: result.Error,
		User:  result.User,
	}
}

func (s *UserStorageService) CreateUser(
	username string,
	password string,
	result chan common.WrappedError,
) {
	userResultChan := make(chan repositories.UserResult)
	go s.repository.GetUser(username, userResultChan)
	userResult := <-userResultChan
	if userResult.Error.Error == nil {
		/// no error means that user already exists
		result <- common.WrappedError{
			Code:  common.AlreadyExists,
			Error: fmt.Errorf("User %s already exists", username),
		}
		return
	} else if userResult.Error.Code != common.NotFound {
		result <- userResult.Error
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		result <- common.WrappedError{
			Code:  common.Unknown,
			Error: err,
		}
	}

	errChan := make(chan common.WrappedError)
	go s.repository.CreateUser(username, string(passwordHash), errChan)
	wrappedErr := <-errChan
	if wrappedErr.Error != nil {
		result <- wrappedErr
	}

	go s.repository.GetUser(username, userResultChan)
	userResult = <-userResultChan
	result <- userResult.Error

}
