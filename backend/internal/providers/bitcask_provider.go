package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"git.mills.io/prologic/bitcask"
)

var bitcaskProviderLock = &sync.Mutex{}
var db *bitcask.Bitcask

func MakeBitcask() *bitcask.Bitcask {
	bitcaskProviderLock.Lock()
	defer bitcaskProviderLock.Unlock()

	if db == nil {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("could not get home directory"))
		}

		dir := filepath.Join(home, ".hyperledger-ui-data", "bitcask")

		db, err = bitcask.Open(dir)
		if err != nil {
			panic(fmt.Sprintf("could not open database at %s because of error: %s", dir, err.Error()))
		}
	}
	return db
}
