package repositories

import (
	"backend/internal/common"
)

type GetConnectionResult struct {
	Data string
	Err  common.WrappedError
}

type IConnectionRepository interface {
	AddConnection(data string, err chan error)
	// returns NotFound if no connection was saved
	GetConnection(oData chan GetConnectionResult)
}
