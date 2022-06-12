package peer_material_services

import "time"

type Node struct {
	Id              string
	NextNodeIds     []string
	PreviousNodeIds []string
}

func MakeNode(
	iNextNodeIds []string,
	iPreviousNodeIds []string,
) Node {
	return Node{
		NextNodeIds:     iNextNodeIds,
		PreviousNodeIds: iPreviousNodeIds,
	}
}

type SignatureOption struct {
	NodeId    string
	Signature string
}

type ReceiveMaterialRequestRequest struct {
	RecipientPublicKey string
	MainNodeId         string
	Nodes              map[string]Node
	TransferTime       time.Time
	SenderPublicKey    string
	SignatureOptions   []SignatureOption
}

func MakeNodeSignatureOption(
	iNodeId string,
	iSignature string,
) SignatureOption {
	return SignatureOption{
		NodeId:    iNodeId,
		Signature: iSignature,
	}
}

func MakeReceiveMaterialRequestRequest(
	iRecipientPublicKey string,
	iMainNodeId string,
	iNodes map[string]Node,
	iTransferTime time.Time,
	iSenderPublicKey string,
	iSignatureOptions []SignatureOption,
) ReceiveMaterialRequestRequest {
	return ReceiveMaterialRequestRequest{
		RecipientPublicKey: iRecipientPublicKey,
		MainNodeId:         iMainNodeId,
		Nodes:              iNodes,
		TransferTime:       iTransferTime,
		SignatureOptions:   iSignatureOptions,
		SenderPublicKey:    iSenderPublicKey,
	}
}
