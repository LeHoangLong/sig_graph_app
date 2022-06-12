package drivers

import (
	"backend/internal/common"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.uber.org/dig"
)

func ProvideSmartContractDriverHl(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(
		func(iConfig common.Config) (SmartContractDriverHL, error) {
			wallet, err := gateway.NewFileSystemWallet(iConfig.WalletPath)
			if err != nil {
				return SmartContractDriverHL{}, err
			}

			identityService := MakeHLIdentityService(nil)

			initializer := MakeHLGatewayInitializer(
				wallet,
				ConfigFilePath(iConfig.HLConfigPath),
				WalletIdentity(iConfig.HLWalletIdentity),
				identityService,
				MspId(iConfig.OrgMspId),
			)
			return MakeSmartContractDriverHL(
				initializer,
				ChannelName(iConfig.ChannelName),
				ContractName(iConfig.ContractName),
			), nil
		},
		dig.As(new(SmartContractDriverI)),
	)
}
