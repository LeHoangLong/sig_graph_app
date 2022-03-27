package models

type PublicKey struct {
	Id    int    `json:"Id"`
	Value string `json:"Value"`
}

func MakePublicKey(
	iId int,
	iValue string,
) PublicKey {
	return PublicKey{
		Id:    iId,
		Value: iValue,
	}
}
