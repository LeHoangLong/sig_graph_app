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

type SenderEndpoint struct {
	Protocol     string
	MajorVersion int
	MinorVersion int
	Url          string
}

type ReceiveMaterialRequestRequest struct {
	RecipientPublicKey string
	MainNodeId         string
	Nodes              map[string]Node
	TransferTime       time.Time
	SenderPublicKey    string
	SignatureOptions   []SignatureOption
	SenderEndpoints    []SenderEndpoint
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

func MakeSenderEndpoint(
	iProtocol string,
	iMajorVersion int,
	iMinorVersion int,
	iUrl string,
) SenderEndpoint {
	return SenderEndpoint{
		Protocol:     iProtocol,
		MajorVersion: iMajorVersion,
		MinorVersion: iMinorVersion,
		Url:          iUrl,
	}
}

func MakeReceiveMaterialRequestRequest(
	iRecipientPublicKey string,
	iMainNodeId string,
	iNodes map[string]Node,
	iTransferTime time.Time,
	iSenderPublicKey string,
	iSignatureOptions []SignatureOption,
	iSenderEndpoints []SenderEndpoint,
) ReceiveMaterialRequestRequest {
	return ReceiveMaterialRequestRequest{
		RecipientPublicKey: iRecipientPublicKey,
		MainNodeId:         iMainNodeId,
		Nodes:              iNodes,
		TransferTime:       iTransferTime,
		SignatureOptions:   iSignatureOptions,
		SenderPublicKey:    iSenderPublicKey,
		SenderEndpoints:    iSenderEndpoints,
	}
}
