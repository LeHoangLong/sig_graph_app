package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
)

type PublicKey string
type PrivateKey string

type GraphContractSignature struct {
	publicKey  PublicKey
	privateKey PrivateKey
}

type NodeSC struct {
	Id                    string          `json:"Id"`
	IsFinalized           bool            `json:"IsFinalized"`
	Data                  interface{}     `json:"Data"`
	NextNodeHashedIds     map[string]bool `json:"NextNodeHashedIds"` /// used a set
	PreviousNodeHashedIds map[string]bool `json:"PreviousNodeHashedIds"`
	OwnerPublicKey        string          `json:"OwnerPublicKey"`
}

func MakeGraphContractSignature(
	iPublicKey PublicKey,
	iPrivateKey PrivateKey,
) GraphContractSignature {
	return GraphContractSignature{
		publicKey:  iPublicKey,
		privateKey: iPrivateKey,
	}
}

func (s GraphContractSignature) CreateNodeSignature(
	iId string,
	iIsFinalized bool,
	iData interface{},
) (string, error) {

	block, _ := pem.Decode([]byte(s.privateKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return "", err
	}

	newNode := NodeSC{
		Id:                    iId,
		Data:                  iData,
		NextNodeHashedIds:     map[string]bool{},
		PreviousNodeHashedIds: map[string]bool{},
		OwnerPublicKey:        string(s.publicKey),
		IsFinalized:           false,
	}

	json, err := json.Marshal(newNode)
	if err != nil {
		return "", err
	}
	hash := sha512.Sum512(json)
	rng := rand.Reader
	signature, err := rsa.SignPKCS1v15(rng, privateKey, crypto.SHA512, hash[:])

	return string(signature), err
}
