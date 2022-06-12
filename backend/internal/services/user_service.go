package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
)

type UserService struct {
	repository repositories.UserRepositoryI
}

func MakeUserService(
	iRepository repositories.UserRepositoryI,
) UserService {
	return UserService{
		repository: iRepository,
	}
}

// func (s *UserService) VerifyUser(username string, password string) error {
// 	user, err := s.repository.GetUser(username)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
// 	return err
// }
//
// func (s *UserService) SignUp(username string, password string) error {
// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = s.repository.CreateUser(username, string(passwordHash))
// 	return err
// }
//
// func (s *UserService) DoesUserExist(username string) (bool, error) {
// 	return s.repository.DoesUserExist(username)
// }

func (s UserService) FindUserWithPublicKey(iPublicKey string) (models.User, error) {
	return s.repository.FindUserWithPublicKey(iPublicKey)
}

func (s UserService) GetUserById(iContext context.Context, iId int) (models.User, error) {
	return s.repository.GetUserById(iContext, iId)
}
