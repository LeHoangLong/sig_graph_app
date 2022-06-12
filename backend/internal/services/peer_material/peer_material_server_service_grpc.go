package peer_material_services

import (
	"backend/internal/grpc"
	"context"
)

type PeerMaterialServerServiceGrpc struct {
	handler ReceiveMaterialRequestReceivedHandler
	grpc.UnimplementedMaterialServiceServer
}

func MakePeerMaterialServerServiceGrpc(
	iHandler ReceiveMaterialRequestReceivedHandler,
) *PeerMaterialServerServiceGrpc {
	return &PeerMaterialServerServiceGrpc{
		handler: iHandler,
	}
}

func convertGrpcToOption(
	iOption *grpc.SignatureOption,
) SignatureOption {
	return SignatureOption{
		NodeId:    iOption.NodeId,
		Signature: string(iOption.Signature),
	}
}

func convertGrpcToNode(
	iNode *grpc.Node,
) Node {
	return Node{
		Id:              iNode.Id,
		NextNodeIds:     iNode.ChildrenNodeIds,
		PreviousNodeIds: iNode.ParentNodeIds,
	}
}

func convertGrpcToReceiveMaterialRequestRequest(
	iRequest *grpc.ReceiveMaterialRequestRequest,
) ReceiveMaterialRequestRequest {
	options := []SignatureOption{}
	for _, option := range iRequest.SenderSignatureOptions {
		options = append(options, convertGrpcToOption(option))
	}

	nodes := map[string]Node{}
	for _, node := range iRequest.Nodes {
		nodes[node.Id] = convertGrpcToNode(node)
	}

	return MakeReceiveMaterialRequestRequest(
		iRequest.RecipientPublicKey,
		iRequest.MainNodeId,
		nodes,
		iRequest.TransferTime.AsTime(),
		iRequest.SenderPublicKey,
		options,
	)
}

func convertResponseToGrpcResponse(
	iReponse ReceiveMaterialRequestResponse,
) grpc.ReceiveMaterialRequestResponse {
	return grpc.ReceiveMaterialRequestResponse{
		ResponseId:          iReponse.ResponseId,
		RequestAcknowledged: iReponse.IsRequestAcknowledged,
	}
}

func (s PeerMaterialServerServiceGrpc) SendReceiveMaterialRequest(
	iCtx context.Context,
	iRequest *grpc.ReceiveMaterialRequestRequest,
) (*grpc.ReceiveMaterialRequestResponse, error) {
	parsedRequest := convertGrpcToReceiveMaterialRequestRequest(iRequest)

	response, err := s.handler.Handle(iCtx, parsedRequest)
	if err != nil {
		return nil, err
	}

	grpcResponse := convertResponseToGrpcResponse(response)
	return &grpcResponse, nil
}
