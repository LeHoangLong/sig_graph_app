package repositories

import (
	"backend/internal/common"
	"errors"

	"git.mills.io/prologic/bitcask"
)

type ConnectionRepositoryBitcask struct {
	db *bitcask.Bitcask
}

func MakeConnectionRepositoryBitcask(
	db *bitcask.Bitcask,
) *ConnectionRepositoryBitcask {
	return &ConnectionRepositoryBitcask{
		db: db,
	}
}

func (c *ConnectionRepositoryBitcask) AddConnection(data string, err chan error) {
	res := c.db.Put([]byte(ConnectionTable), []byte(data))
	err <- res
}

func (c *ConnectionRepositoryBitcask) GetConnection(oData chan GetConnectionResult) {
	dataBytes, err := c.db.Get([]byte(ConnectionTable))
	data := GetConnectionResult{}
	if errors.Is(err, bitcask.ErrKeyNotFound) {
		data.Data = ""
		data.Err = common.WrappedError{
			Code: common.NotFound,
		}
	} else {
		str := string(dataBytes)
		data.Data = str
		data.Err = common.WrappedError{
			Error: err,
		}
	}
	oData <- data
}
