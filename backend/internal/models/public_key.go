package models

type PublicKeyId int

/// Id may be nil if it is obtained from blockchain
type PublicKey struct {
	Id    *PublicKeyId `json:"Id"`
	Value string       `json:"Value"`
}

func MakePublicKey(
	iId *PublicKeyId,
	iValue string,
) PublicKey {
	return PublicKey{
		Id:    iId,
		Value: iValue,
	}
}
