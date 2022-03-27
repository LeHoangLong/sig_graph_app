package drivers

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type SmartContractDriverHL struct {
	contract *gateway.Contract
	gateway  *gateway.Gateway /// gateway assocaited with contract
	mtx      *sync.Mutex
}

type ChannelName string
type ContractName string

func MakeSmartContractDriverHL(
	iInitializer HLGatewayInitializer,
	iChannelName ChannelName,
	iContractName ContractName,
) SmartContractDriverHL {
	gateway, err := iInitializer.MakeHLGateway()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize gateway %s", err.Error()))
	}

	network, err := gateway.GetNetwork("mychannel")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize gateway network %s", err.Error()))
	}
	contract := network.GetContractWithName(string(iContractName), "MaterialContract")

	driver := SmartContractDriverHL{
		contract,
		gateway,
		&sync.Mutex{},
	}

	runtime.SetFinalizer(&driver, close)
	return driver
}

func close(d *SmartContractDriverHL) {
	d.gateway.Close()
}

func (d SmartContractDriverHL) CreateTransaction(
	iFunctionName string,
	iArgs ...string,
) ([]byte, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	tx, err := d.contract.CreateTransaction(iFunctionName)
	if err != nil {
		return []byte{}, err
	}
	return tx.Submit(iArgs...)
}

func (d SmartContractDriverHL) Query(
	iFunctionName string,
	iArgs ...string,
) ([]byte, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	return d.contract.EvaluateTransaction(iFunctionName, iArgs...)
}
