package controllers

import (
	"backend/internal/providers"
	"backend/internal/services"
)

type ConnectionController struct {
	hyperledgerConnectionService *services.ConnectHyperledgerService
	connectionStorageService     *services.ConnectionStorageService
	userStorageService           *services.UserStorageService
}

func MakeController(
	connectionStorageService *services.ConnectionStorageService,
	userStorageService *services.UserStorageService,
) *ConnectionController {
	return &ConnectionController{
		hyperledgerConnectionService: nil,
		connectionStorageService:     connectionStorageService,
		userStorageService:           userStorageService,
	}
}

func (c *ConnectionController) SaveConnection(path string) (string, error) {
	errChan := make(chan error)
	go c.connectionStorageService.AddConnection(path, errChan)
	err := <-errChan
	if err != nil {
		return "", err
	}

	return c.LoadConnection()
}

func (c *ConnectionController) LoadConnection() (string, error) {
	connectionResultChan := make(chan services.GetConnectionResult)
	go c.connectionStorageService.GetConnection(connectionResultChan)
	connectionResult := <-connectionResultChan
	return connectionResult.Data, connectionResult.Err.Error
}

func (c *ConnectionController) Connect() error {
	configStr, err := c.LoadConnection()
	if err != nil {
		return err
	}

	if c.hyperledgerConnectionService == nil {
		userResultChan := make(chan services.UserResult)
		go c.userStorageService.GetCurrentUser(userResultChan)
		userResult := <-userResultChan
		if userResult.Error.Error != nil {
			return userResult.Error.Error
		}

		wallet := providers.MakeHyperledgerWallet()
		c.hyperledgerConnectionService = services.MakeConnectHyperledgerService(
			wallet,
			userResult.User.Username,
		)
	}

	return c.hyperledgerConnectionService.Connect(configStr)
}
