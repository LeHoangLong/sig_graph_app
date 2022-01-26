package controllers

import (
	"backend/internal/common"
	"backend/internal/services"
)

type UserController struct {
	storageService    *services.UserStorageService
	enrollmentService *services.EnrollUserHyperledgerSerivce
}

func MakeUserController(
	storageService *services.UserStorageService,
	enrollmentService *services.EnrollUserHyperledgerSerivce,
) *UserController {
	return &UserController{
		storageService:    storageService,
		enrollmentService: enrollmentService,
	}
}

func (c *UserController) CreateUser(
	username string,
	password string,
	organizationMspId string,
	publicCertPath string,
	privateKeyPath string,
) error {
	resultChan := make(chan common.WrappedError)
	go c.storageService.CreateUser(username, password, resultChan)
	result := <-resultChan
	if result.Error != nil {
		return result.Error
	}

	go c.enrollmentService.CreateIdentity(
		username,
		organizationMspId,
		publicCertPath,
		privateKeyPath,
		resultChan,
	)

	result = <-resultChan
	return result.Error
}

func (c *UserController) Login(username string, password string) error {
	errorChan := make(chan common.WrappedError)
	go c.storageService.VerifyUser(username, password, errorChan)
	err := (<-errorChan).Error
	if err != nil {
		return err
	}

	go c.storageService.SetCurrentUser(username, errorChan)
	return (<-errorChan).Error
}
