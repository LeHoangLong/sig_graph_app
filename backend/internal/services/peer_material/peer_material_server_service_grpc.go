package peer_material_services

import (
	"backend/internal/grpc"
	"context"
)

type PeerMaterialServerServiceGrpc struct {
	handler         ReceiveMaterialRequestReceivedHandler
	responseHandler ReceiveMaterialResponseHandler
	grpc.UnimplementedMaterialServiceServer
}

func MakePeerMaterialServerServiceGrpc(
	iHandler ReceiveMaterialRequestReceivedHandler,
	iResponseHandler ReceiveMaterialResponseHandler,
) *PeerMaterialServerServiceGrpc {
	return &PeerMaterialServerServiceGrpc{
		handler:         iHandler,
		responseHandler: iResponseHandler,
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

func convertGrpcToSenderEndpoint(
	iSenderEndpoint *grpc.Endpoint,
) SenderEndpoint {
	return MakeSenderEndpoint(
		iSenderEndpoint.Protocol,
		int(iSenderEndpoint.MajorVersion),
		int(iSenderEndpoint.MinorVersion),
		iSenderEndpoint.Url,
	)
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

	senderEndpoints := []SenderEndpoint{}
	for _, endpoint := range iRequest.SenderEndpoints {
		senderEndpoints = append(senderEndpoints, convertGrpcToSenderEndpoint(endpoint))
	}

	return MakeReceiveMaterialRequestRequest(
		iRequest.RecipientPublicKey,
		iRequest.MainNodeId,
		nodes,
		iRequest.TransferTime.AsTime(),
		iRequest.SenderPublicKey,
		options,
		senderEndpoints,
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

func convertGrpcToReceiveMaterialResponse(
	iResponse *grpc.ReceiveMaterialResponseRequest,
) ReceiveMaterialResponse {
	response := MakeReceiveMaterialResponse(
		iResponse.ResponseId,
		iResponse.IsRequestAccepted,
		iResponse.Message,
		iResponse.NewNodeId,
	)
	return response
}

func makeGrpcReceiveMaterialResponseAcknowledgement() grpc.ReceiveMaterialResponseResponse {
	return grpc.ReceiveMaterialResponseResponse{}
}

func convertReceiveMaterialResponseAcknowledgementToGrpc(
	iAck ReceiveMaterialResponseAcknowledgement,
) grpc.ReceiveMaterialResponseResponse {
	return makeGrpcReceiveMaterialResponseAcknowledgement()
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

func (s PeerMaterialServerServiceGrpc) SendReceiveMaterialResponse(
	ctx context.Context,
	in *grpc.ReceiveMaterialResponseRequest,
) (*grpc.ReceiveMaterialResponseResponse, error) {
	response := convertGrpcToReceiveMaterialResponse(in)
	ack, err := s.responseHandler.HandleReponse(
		ctx,
		response,
	)
	if err != nil {
		return nil, err
	}

	parsedAck := convertReceiveMaterialResponseAcknowledgementToGrpc(ack)
	return &parsedAck, nil
}
