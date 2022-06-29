package models

type UserEndpointId int
type UserEndpoint struct {
	Id       UserEndpointId
	UserId   UserId
	Url      string
	Protocol PeerProtocol
}

func MakeUserEndpoint(
	iId UserEndpointId,
	iUserId UserId,
	iUrl string,
	iProtocol PeerProtocol,
) UserEndpoint {
	return UserEndpoint{
		Id:       iId,
		UserId:   iUserId,
		Url:      iUrl,
		Protocol: iProtocol,
	}
}
