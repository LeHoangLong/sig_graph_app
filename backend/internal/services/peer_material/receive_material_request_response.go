package peer_material_services

type ReceiveMaterialRequestResponse struct {
	ResponseId            string
	IsRequestAcknowledged bool
}

func MakeReceiveMaterialRequestResponse(
	iResponseId string,
	iIsRequestAcknowledged bool,
) ReceiveMaterialRequestResponse {
	return ReceiveMaterialRequestResponse{
		ResponseId:            iResponseId,
		IsRequestAcknowledged: iIsRequestAcknowledged,
	}
}
