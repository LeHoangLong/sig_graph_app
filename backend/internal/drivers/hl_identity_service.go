package drivers

import (
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type HLIdentityService struct {
	client *msp.Client
}

type MspId string
type Username string

func MakeHLIdentityService(
	iClient *msp.Client,
) HLIdentityService {
	return HLIdentityService{
		client: iClient,
	}
}

func (d HLIdentityService) CreateX509CertificateFromFiles(
	iOrgMspId MspId,
	iUsername Username,
) (*gateway.X509Identity, error) {
	credPath := filepath.Join(
		"keystore",
		string(iUsername),
	)
	certPath := filepath.Join(credPath, "public.pem")
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, nil
	}

	keyPath := filepath.Join(credPath, "private.pem")
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return nil, nil
	}

	return gateway.NewX509Identity(string(iOrgMspId), string(cert), string(key)), nil
}

func (d HLIdentityService) createX509Certificate(
	iUsername Username,
) (*gateway.X509Identity, error) {
	secret, err := d.client.Register(
		&msp.RegistrationRequest{
			Name: string(iUsername),
		},
	)

	if err != nil {
		return nil, err
	}

	err = d.client.Enroll(
		string(iUsername),
		msp.WithSecret(secret),
	)

	if err != nil {
		return nil, err
	}

	identity, err := d.client.GetSigningIdentity(string(iUsername))
	if err != nil {
		return nil, err
	}

	privateKey, err := identity.PrivateKey().Bytes()
	if err != nil {
		return nil, err
	}

	x509Identity := gateway.NewX509Identity(
		identity.PublicVersion().Identifier().MSPID,
		string(identity.PublicVersion().EnrollmentCertificate()),
		string(privateKey),
	)

	return x509Identity, nil
}
