package models

type ProtocolName string

const (
	GRPC ProtocolName = "grpc"
)

type PeerProtocolId int
type PeerProtocol struct {
	Id           PeerProtocolId
	Name         ProtocolName
	MajorVersion int
	MinorVersion int
}

func MakePeerProtocol(
	iId PeerProtocolId,
	iProtocol ProtocolName,
	iMajorVersion int,
	iMinorVersion int,
) PeerProtocol {
	return PeerProtocol{
		Id:           iId,
		Name:         iProtocol,
		MajorVersion: iMajorVersion,
		MinorVersion: iMinorVersion,
	}
}
