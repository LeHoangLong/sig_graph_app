package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var wallet *gateway.Wallet
var hyperledgerWalletProviderLock = &sync.Mutex{}

func MakeHyperledgerWallet() *gateway.Wallet {
	hyperledgerWalletProviderLock.Lock()
	defer hyperledgerWalletProviderLock.Unlock()

	if wallet == nil {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("could not get home directory"))
		}

		dir := filepath.Join(home, ".hyperledger-wallet")
		wallet, err = gateway.NewFileSystemWallet(dir)
		if err != nil {
			panic(fmt.Sprintf("could not open wallet at %s because of error: %s", dir, err.Error()))
		}
	}
	return wallet
}
