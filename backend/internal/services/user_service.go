package services

import (
	"backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.UserRepositoryI
}

func (s *UserService) VerifyUser(username string, password string) error {
	user, err := s.repository.GetUser(username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err
}

func (s *UserService) SignUp(username string, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	_, err = s.repository.CreateUser(username, string(passwordHash))
	return err
}

func (s *UserService) DoesUserExist(username string) error {
	return s.DoesUserExist(username)
}
