package node_contract

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
)

type SignatureService struct {
}

func MakeSignatureService() SignatureService {
	return SignatureService{}
}

/// iNode's signature will be ignored
func (s SignatureService) CreateNodeSignature(
	iPrivatekey string,
	iNode NodeI,
) (string, error) {
	block, _ := pem.Decode([]byte(iPrivatekey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return "", err
	}

	noSignatureHeader := iNode.GetHeader()
	tempHeader := iNode.GetHeader()
	noSignatureHeader.Signature = ""
	defer func() {
		iNode.SetHeader(tempHeader)
	}()
	iNode.SetHeader(noSignatureHeader)
	nodeJson, err := json.Marshal(iNode)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(nodeJson)
	rng := rand.Reader
	signature, err := rsa.SignPKCS1v15(rng, privateKey, crypto.SHA512, hash[:])

	return string(signature), err
}
