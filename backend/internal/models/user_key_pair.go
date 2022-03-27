package models

type UserKeyPair struct {
	Id         int       `json:"Id"`
	PublicKey  PublicKey `json:"PublicKey"`
	PrivateKey string    `json:"PrivateKey"`
	IsDefault  bool      `json:"IsDefault"`
}

func MakeUserKeyPair(
	iId int,
	iPublicKey PublicKey,
	iPrivateKey string,
	iIsDefault bool,
) UserKeyPair {
	return UserKeyPair{
		Id:         iId,
		PublicKey:  iPublicKey,
		PrivateKey: iPrivateKey,
		IsDefault:  iIsDefault,
	}
}
