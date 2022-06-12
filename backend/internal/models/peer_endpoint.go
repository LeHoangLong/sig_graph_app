package models

type PeerEndpoint struct {
	Id       int
	Url      string
	Protocol PeerProtocol
}

func MakePeerEndpoint(
	iId int,
	iUrl string,
	iProtocol PeerProtocol,
) PeerEndpoint {
	return PeerEndpoint{
		Id:       iId,
		Url:      iUrl,
		Protocol: iProtocol,
	}
}
