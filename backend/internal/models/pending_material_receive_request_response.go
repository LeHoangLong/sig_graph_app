package models

type PendingMaterialReceiveRequestResponse struct {
	RequestId  int
	ResponseId string
}

func MakePendingMaterialReceiveRequestResponse(
	iRequestId int,
	iResponseId string,
) PendingMaterialReceiveRequestResponse {
	return PendingMaterialReceiveRequestResponse{
		RequestId:  iRequestId,
		ResponseId: iResponseId,
	}
}
