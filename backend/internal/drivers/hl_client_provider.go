package drivers

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type HLClientProvider struct {
	sdk *fabsdk.FabricSDK
}

func MakeHLChannelProvider(
	iConfigFilePath ConfigFilePath,
) (HLClientProvider, error) {
	fabSdk, err := fabsdk.New(config.FromFile(string(iConfigFilePath)))
	if err != nil {
		return HLClientProvider{}, err
	}
	return HLClientProvider{
		sdk: fabSdk,
	}, nil
}

func (d HLClientProvider) GetClient() (*msp.Client, error) {
	provider := d.sdk.Context()
	client, err := msp.New(provider)
	if err != nil {
		return nil, err
	}

	return client, err
}
