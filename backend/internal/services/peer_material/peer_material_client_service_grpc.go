package peer_material_services

import (
	"backend/internal/grpc"
	"context"

	google_grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PeerMaterialClientServiceGrpc struct {
	client grpc.MaterialServiceClient
}

func MakePeerMaterialClientServiceGrpc(
	iConnection google_grpc.ClientConnInterface,
) PeerMaterialClientServiceGrpc {
	return PeerMaterialClientServiceGrpc{
		client: grpc.NewMaterialServiceClient(iConnection),
	}
}

func convertOptionToGrpcOption(
	iOption SignatureOption,
) *grpc.SignatureOption {
	return &grpc.SignatureOption{
		NodeId:    iOption.NodeId,
		Signature: []byte(iOption.Signature),
	}
}

func convertNodeToGrpc(
	iNode Node,
) *grpc.Node {
	return &grpc.Node{
		Id:              iNode.Id,
		ChildrenNodeIds: iNode.NextNodeIds,
		ParentNodeIds:   iNode.PreviousNodeIds,
	}
}

func convertReceiverMaterialRequestToGrpc(
	iRequest ReceiveMaterialRequestRequest,
) *grpc.ReceiveMaterialRequestRequest {
	grpcOptions := []*grpc.SignatureOption{}
	for _, option := range iRequest.SignatureOptions {
		grpcOptions = append(grpcOptions, convertOptionToGrpcOption(option))
	}

	grpcNodes := map[string]*grpc.Node{}
	for _, node := range iRequest.Nodes {
		grpcNodes[node.Id] = convertNodeToGrpc(node)
	}

	return &grpc.ReceiveMaterialRequestRequest{
		RecipientPublicKey:     iRequest.RecipientPublicKey,
		MainNodeId:             iRequest.MainNodeId,
		TransferTime:           timestamppb.New(iRequest.TransferTime),
		SenderSignatureOptions: grpcOptions,
		SenderPublicKey:        iRequest.SenderPublicKey,
		Nodes:                  grpcNodes,
	}
}

func (s PeerMaterialClientServiceGrpc) SendReceiveMaterialRequest(
	iCtx context.Context,
	iRequest ReceiveMaterialRequestRequest,
) (ReceiveMaterialRequestResponse, error) {
	grpcOptions := []*grpc.SignatureOption{}
	for _, option := range iRequest.SignatureOptions {
		grpcOptions = append(grpcOptions, convertOptionToGrpcOption(option))
	}

	request := convertReceiverMaterialRequestToGrpc(iRequest)
	response, err := s.client.SendReceiveMaterialRequest(iCtx, request)
	if err != nil {
		return ReceiveMaterialRequestResponse{}, err
	}

	return MakeReceiveMaterialRequestResponse(response.ResponseId, response.RequestAcknowledged), nil
}
