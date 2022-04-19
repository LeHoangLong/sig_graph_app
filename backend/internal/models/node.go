package models

type Node struct {
	Id                    *int            `json:"Id"`
	NodeId                string          `json:"NodeId"`
	IsFinalized           bool            `json:"IsFinalized"`
	PreviousNodeHashedIds map[string]bool `json:"PreviousNodeHashedIds"` /// used as a set
	NextNodeHashedIds     map[string]bool `json:"NextodeHashedIds"`      /// used as a set
	OwnerPublicKey        PublicKey       `json:"OwnerPublicKey"`
	CreatedTime           CustomTime      `json:"CreatedTime"`
	Signature             string          `json:"Signature"`
}

func MakeNode(
	iId *int,
	iNodeId string,
	iIsFinalized bool,
	iPreviousNodeHashedIds map[string]bool,
	iNextNodeHashedIds map[string]bool,
	iOwnerPublicKey PublicKey,
	iCreatedTime CustomTime,
	iSignature string,
) Node {
	return Node{
		Id:                    iId,
		NodeId:                iNodeId,
		IsFinalized:           iIsFinalized,
		PreviousNodeHashedIds: iPreviousNodeHashedIds,
		NextNodeHashedIds:     iNextNodeHashedIds,
		OwnerPublicKey:        iOwnerPublicKey,
		CreatedTime:           iCreatedTime,
		Signature:             iSignature,
	}
}
