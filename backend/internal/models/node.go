package models

type Node struct {
	Id                    *int            `json:"Id"`
	NodeId                string          `json:"NodeId"`
	IsFinalized           bool            `json:"IsFinalized"`
	PreviousNodeHashedIds map[string]bool `json:"PreviousNodeHashedIds"` /// used as a set
	NextNodeHashedIds     map[string]bool `json:"NextNodeHashedIds"`
	ChildrenIds           map[int]bool
	ParentIds             map[int]bool
	OwnerPublicKey        PublicKey  `json:"OwnerPublicKey"`
	CreatedTime           CustomTime `json:"CreatedTime"`
	Signature             string     `json:"Signature"`
	Type                  string     `json:"Type"`
}

func MakeNode(
	iId *int,
	iNodeId string,
	iIsFinalized bool,
	iPreviousNodeHashedIds map[string]bool,
	iNextNodeHashedIds map[string]bool,
	iParentIds map[int]bool,
	iChildrenIds map[int]bool,
	iOwnerPublicKey PublicKey,
	iCreatedTime CustomTime,
	iSignature string,
	iType string,
) Node {
	return Node{
		Id:                    iId,
		NodeId:                iNodeId,
		IsFinalized:           iIsFinalized,
		PreviousNodeHashedIds: iPreviousNodeHashedIds,
		NextNodeHashedIds:     iNextNodeHashedIds,
		ParentIds:             iParentIds,
		ChildrenIds:           iChildrenIds,
		OwnerPublicKey:        iOwnerPublicKey,
		CreatedTime:           iCreatedTime,
		Signature:             iSignature,
		Type:                  iType,
	}
}
