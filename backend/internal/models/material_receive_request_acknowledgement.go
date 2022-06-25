package models

type MaterialReceiveRequestResponseId string

type MaterialReceiveRequestAcknowledgement struct {
	RequestId  MaterialReceiveRequestId
	ResponseId MaterialReceiveRequestResponseId
}

func MakeMaterialReceiveRequestAcknowledgement(
	iRequestId MaterialReceiveRequestId,
	iResponseId MaterialReceiveRequestResponseId,
) MaterialReceiveRequestAcknowledgement {
	return MaterialReceiveRequestAcknowledgement{
		RequestId:  iRequestId,
		ResponseId: iResponseId,
	}
}
