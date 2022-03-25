package drivers

import (
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type ConfigFilePath string
type WalletIdentity string

type HLGatewayInitializer struct {
	wallet          *gateway.Wallet
	configFilePath  ConfigFilePath
	walletIdentity  WalletIdentity
	identityService HLIdentityService
	orgMspId        MspId
}

func MakeHLGatewayInitializer(
	iWallet *gateway.Wallet,
	iConfigFilePath ConfigFilePath,
	iWalletIdentity WalletIdentity,
	iHlIdentityService HLIdentityService,
	iOrgMspId MspId,
) HLGatewayInitializer {
	return HLGatewayInitializer{
		wallet:          iWallet,
		configFilePath:  iConfigFilePath,
		walletIdentity:  iWalletIdentity,
		identityService: iHlIdentityService,
		orgMspId:        iOrgMspId,
	}
}

func (s HLGatewayInitializer) MakeHLGateway() (*gateway.Gateway, error) {
	if !s.wallet.Exists(string(s.walletIdentity)) {
		identity, err := s.identityService.CreateX509CertificateFromFiles(
			s.orgMspId,
			Username(s.walletIdentity),
		)
		if err != nil {
			return nil, err
		}
		s.wallet.Put(string(s.walletIdentity), identity)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(string(s.configFilePath)))),
		gateway.WithIdentity(s.wallet, string(s.walletIdentity)),
	)
	if err != nil {
		return nil, err
	}

	return gw, err
}
