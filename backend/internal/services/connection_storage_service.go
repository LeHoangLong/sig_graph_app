package services

import (
	"backend/internal/common"
	"backend/internal/repositories"
	"os"
)

type ConnectionStorageService struct {
	repository           repositories.IConnectionRepository
	predefinedConfigPath string
}

func MakeConnectionStorageService(
	repository repositories.IConnectionRepository,
	predefinedConfigPath string,
) *ConnectionStorageService {
	return &ConnectionStorageService{
		repository:           repository,
		predefinedConfigPath: predefinedConfigPath,
	}
}

func (c *ConnectionStorageService) AddConnection(path string, retError chan error) {
	data, err := os.ReadFile(path)
	if err == nil {
		err := make(chan error)
		go c.repository.AddConnection(string(data), err)
		retError <- <-err
	} else {
		retError <- err
	}
}

type GetConnectionResult struct {
	Data string
	Err  common.WrappedError
}

func (c *ConnectionStorageService) GetConnection(oData chan GetConnectionResult) {
	repoResultChan := make(chan repositories.GetConnectionResult)
	go c.repository.GetConnection(repoResultChan)
	repoResult := <-repoResultChan

	ret := GetConnectionResult{}
	if repoResult.Err.Code == common.NotFound {
		ret.Data = ""
		ret.Err = common.WrappedError{}
	} else {
		ret.Data = repoResult.Data
		ret.Err = repoResult.Err
	}

	oData <- ret
}

func (c *ConnectionStorageService) GetPredefinedConnection(oData chan GetConnectionResult) {
	data, err := os.ReadFile(c.predefinedConfigPath)
	if err != nil {
		oData <- GetConnectionResult{
			Err: common.WrappedError{
				Error: err,
			},
		}
		return
	}

	oData <- GetConnectionResult{
		Data: string(data),
	}
}
