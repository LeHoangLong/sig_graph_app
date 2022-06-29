package repositories

import (
	"backend/internal/models"
	"context"
)

type ProtocolWithMajorVersion struct {
	Protocol     string
	MajorVersion int
}

func MakeProtocolWithMajorVersion(
	iProtocol string,
	iMajorVersion int,
) ProtocolWithMajorVersion {
	return ProtocolWithMajorVersion{
		Protocol:     iProtocol,
		MajorVersion: iMajorVersion,
	}
}

type PeerProtocolRepositoryI interface {
	FetchPeerProtocolByIds(
		iContext context.Context,
		iIds map[models.PeerProtocolId]bool,
	) (map[models.PeerProtocolId]models.PeerProtocol, error)
	/// each protocol with major version may have multiple
	/// minor versions
	FilterSupportedProtocolsWithMajorVersion(
		iContext context.Context,
		iProtocols map[ProtocolWithMajorVersion]bool,
	) (map[ProtocolWithMajorVersion][]models.PeerProtocol, error)
}
