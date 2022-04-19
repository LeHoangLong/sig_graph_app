package node_contract

import (
	"backend/internal/models"
)

type NodeHeader struct {
	Id                    string            `json:"Id"`
	IsFinalized           bool              `json:"IsFinalized"`
	PreviousNodeHashedIds map[string]bool   `json:"PreviousNodeHashedIds"` /// used as a set
	NextNodeHashedIds     map[string]bool   `json:"NextodeHashedIds"`      /// used as a set
	OwnerPublicKey        string            `json:"OwnerPublicKey"`
	CreatedTime           models.CustomTime `json:"CreatedTime"`
	Signature             string            `json:"Signature"`
}

type NodeI interface {
	GetHeader() NodeHeader
	SetHeader(NodeHeader)
}

func MakeNodeHeader(
	iId string,
	iIsFinalized bool,
	iPreviousNodeHashedIds map[string]bool,
	iNextNodeHashedIds map[string]bool,
	iOwnerPublicKey string,
	iCreatedTime models.CustomTime,
	iSignature string,
) NodeHeader {
	return NodeHeader{
		Id:                    iId,
		IsFinalized:           iIsFinalized,
		NextNodeHashedIds:     iNextNodeHashedIds,
		PreviousNodeHashedIds: iPreviousNodeHashedIds,
		OwnerPublicKey:        iOwnerPublicKey,
		CreatedTime:           iCreatedTime,
		Signature:             iSignature,
	}
}
