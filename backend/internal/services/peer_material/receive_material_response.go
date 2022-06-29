package peer_material_services

type ReceiveMaterialResponse struct {
	ResponseId        string
	IsRequestAccepted bool
	Message           string
	NewNodeId         string
}

func MakeReceiveMaterialResponse(
	iResponseId string,
	iIsRequestAccepted bool,
	iMessage string,
	iNewNodeId string,
) ReceiveMaterialResponse {
	return ReceiveMaterialResponse{
		ResponseId:        iResponseId,
		IsRequestAccepted: iIsRequestAccepted,
		Message:           iMessage,
		NewNodeId:         iNewNodeId,
	}
}
