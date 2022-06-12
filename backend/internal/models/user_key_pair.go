package models

type UserKeyPair struct {
	PublicKey
	PrivateKey string `json:"PrivateKey"`
	IsDefault  bool   `json:"IsDefault"`
}

func MakeUserKeyPair(
	iPublicKey PublicKey,
	iPrivateKey string,
	iIsDefault bool,
) UserKeyPair {
	return UserKeyPair{
		PublicKey:  iPublicKey,
		PrivateKey: iPrivateKey,
		IsDefault:  iIsDefault,
	}
}
