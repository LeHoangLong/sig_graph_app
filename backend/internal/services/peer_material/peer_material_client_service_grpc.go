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

func makeGrpcEndpoint(
	iProtocol string,
	iMajorVersion int32,
	iMinorVersion int32,
	iUrl string,
) grpc.Endpoint {
	return grpc.Endpoint{
		Protocol:     iProtocol,
		MajorVersion: iMajorVersion,
		MinorVersion: iMinorVersion,
		Url:          iUrl,
	}
}

func convertSenderEndpointsToGrpc(
	iEndpoint SenderEndpoint,
) *grpc.Endpoint {
	endpoint := makeGrpcEndpoint(
		iEndpoint.Protocol,
		int32(iEndpoint.MajorVersion),
		int32(iEndpoint.MinorVersion),
		iEndpoint.Url,
	)
	return &endpoint
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

	endpoints := []*grpc.Endpoint{}
	for i := range iRequest.SenderEndpoints {
		endpoints = append(endpoints, convertSenderEndpointsToGrpc(iRequest.SenderEndpoints[i]))
	}

	return &grpc.ReceiveMaterialRequestRequest{
		RecipientPublicKey:     iRequest.RecipientPublicKey,
		MainNodeId:             iRequest.MainNodeId,
		TransferTime:           timestamppb.New(iRequest.TransferTime),
		SenderSignatureOptions: grpcOptions,
		SenderPublicKey:        iRequest.SenderPublicKey,
		Nodes:                  grpcNodes,
		SenderEndpoints:        endpoints,
	}
}

func makeGrpcReceiveMaterialResponseRequest(
	iResponseId string,
	iIsRequestAccepted bool,
	iMessage string,
	iNewNodeId string,
) grpc.ReceiveMaterialResponseRequest {
	return grpc.ReceiveMaterialResponseRequest{
		ResponseId:        iResponseId,
		IsRequestAccepted: iIsRequestAccepted,
		Message:           iMessage,
		NewNodeId:         iNewNodeId,
	}
}

func convertReceiveMaterialResponseToGrpc(
	iResponse ReceiveMaterialResponse,
) *grpc.ReceiveMaterialResponseRequest {
	request := makeGrpcReceiveMaterialResponseRequest(
		iResponse.ResponseId,
		iResponse.IsRequestAccepted,
		iResponse.Message,
		iResponse.NewNodeId,
	)
	return &request
}

func convertGrpcToReceiveMaterialResponseAcknowledgement(
	iResponse *grpc.ReceiveMaterialResponseResponse,
) ReceiveMaterialResponseAcknowledgement {
	ack := MakeReceiveMaterialResponseAcknowledgement()
	return ack
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

func (s PeerMaterialClientServiceGrpc) SendReceiveMaterialResponse(
	iCtx context.Context,
	iResponse ReceiveMaterialResponse,
) (ReceiveMaterialResponseAcknowledgement, error) {
	response := convertReceiveMaterialResponseToGrpc(iResponse)
	ack, err := s.client.SendReceiveMaterialResponse(iCtx, response)
	if err != nil {
		return ReceiveMaterialResponseAcknowledgement{}, err
	}

	parsedAck := convertGrpcToReceiveMaterialResponseAcknowledgement(ack)
	return parsedAck, nil
}
