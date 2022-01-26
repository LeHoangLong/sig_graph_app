package services

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type ConnectHyperledgerService struct {
	gateway  *gateway.Gateway
	wallet   *gateway.Wallet
	username string
}

func MakeConnectHyperledgerService(
	wallet *gateway.Wallet,
	username string,
) *ConnectHyperledgerService {
	c := &ConnectHyperledgerService{
		wallet:   wallet,
		username: username,
	}

	return c
}

func (c *ConnectHyperledgerService) Connect(
	iConfig string,
) error {
	var err error

	c.gateway, err = gateway.Connect(
		gateway.WithConfig(
			config.FromRaw(
				[]byte(iConfig),
				"json",
			),
		),
		gateway.WithIdentity(c.wallet, c.username),
	)

	return err
}
