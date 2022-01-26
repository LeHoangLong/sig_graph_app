package services

import (
	"backend/internal/common"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type EnrollUserHyperledgerSerivce struct {
	wallet *gateway.Wallet
}

func MakeEnrollUserHyperledgerService(
	wallet *gateway.Wallet,
) *EnrollUserHyperledgerSerivce {
	s := &EnrollUserHyperledgerSerivce{
		wallet: wallet,
	}
	return s
}

/// overwrite if username already exists
func (s *EnrollUserHyperledgerSerivce) CreateIdentity(
	username string,
	organizationMpsId string,
	publicCertificatePath string,
	privateKeyPath string,
	oErr chan common.WrappedError,
) {
	publicCertificate, err := os.ReadFile(publicCertificatePath)
	if err != nil {
		oErr <- common.WrappedError{
			Error: err,
		}
		return
	}

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		oErr <- common.WrappedError{
			Error: err,
		}
		return
	}

	err = s.wallet.Put(
		username,
		gateway.NewX509Identity(
			organizationMpsId,
			string(publicCertificate),
			string(privateKey),
		),
	)

	oErr <- common.WrappedError{
		Error: err,
	}
}
