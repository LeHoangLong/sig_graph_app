package repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"

	"git.mills.io/prologic/bitcask"
)

type UserRepositoryBitcask struct {
	db *bitcask.Bitcask
}

func MakeUserRepositoryBitcask(
	db *bitcask.Bitcask,
) *UserRepositoryBitcask {
	return &UserRepositoryBitcask{
		db: db,
	}
}
func (r *UserRepositoryBitcask) CreateUser(username string, passwordHash string, done chan common.WrappedError) {
	user := models.User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	json, err := json.Marshal(user)
	if err != nil {
		done <- common.WrappedError{
			Code:  common.FailToUnmarshalJson,
			Error: err,
		}
	}
	err = r.db.Put([]byte(fmt.Sprintf("%s/%s", UserTable, username)), json)
	done <- common.WrappedError{
		Code:  common.Unknown,
		Error: err,
	}
}

func (r *UserRepositoryBitcask) SetCurrentUsername(username string, done chan common.WrappedError) {
	err := r.db.Put([]byte(CurrentUserTable), []byte(username))
	done <- common.WrappedError{
		Error: err,
	}
}

func (r *UserRepositoryBitcask) GetCurrentUsername(result chan UsernameResult) {
	data, err := r.db.Get([]byte(CurrentUserTable))
	if err != nil {
		code := common.Unknown
		if errors.Is(err, bitcask.ErrKeyNotFound) {
			code = common.NotFound
		}
		result <- UsernameResult{
			Username: "",
			Error: common.WrappedError{
				Code:  code,
				Error: err,
			},
		}
		return
	}

	result <- UsernameResult{
		Username: string(data),
		Error: common.WrappedError{
			Error: nil,
		},
	}
}

func (r *UserRepositoryBitcask) GetUser(username string, oUser chan UserResult) {
	data, err := r.db.Get([]byte(fmt.Sprintf("%s/%s", UserTable, username)))
	ret := UserResult{}
	if errors.Is(err, bitcask.ErrKeyNotFound) {
		ret.User = nil
		ret.Error = common.WrappedError{
			Code:  common.NotFound,
			Error: err,
		}
		oUser <- ret
		return
	} else if err != nil {
		ret.User = nil
		ret.Error = common.WrappedError{
			Code:  common.Unknown,
			Error: err,
		}
		oUser <- ret
		return
	}

	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		ret.User = nil
		ret.Error = common.WrappedError{
			Code:  common.FailToUnmarshalJson,
			Error: err,
		}

		oUser <- ret
		return
	}

	ret.User = &user
	ret.Error = common.WrappedError{}

	oUser <- ret
}
