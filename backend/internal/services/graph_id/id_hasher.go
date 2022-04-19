package graph_id_service

import (
	"crypto/sha512"
	"hash"
)

type IdHasher struct {
	hasher hash.Hash
}

func MakeIdHasher() IdHasher {
	return IdHasher{
		hasher: sha512.New(),
	}
}

func (h IdHasher) HashId(
	iId string,
) []byte {
	return h.hasher.Sum([]byte(iId))
}
