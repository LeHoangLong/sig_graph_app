package node_contract

import "crypto/sha512"

type IdHasher struct {
}

func MakeIdHasher() IdHasher {
	return IdHasher{}
}

func (h IdHasher) Hash(iId string) string {
	temp := sha512.Sum512([]byte(iId))
	return string(temp[:])
}
