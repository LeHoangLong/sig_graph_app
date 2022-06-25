package repositories

import (
	"backend/internal/models"
	"context"
)

type PeerRepositoryI interface {
	FetchPeers(iUserId models.UserId) ([]models.Peer, error)
	FetchPeerEndPoints(iContext context.Context, iPeerId models.PeerId) ([]models.PeerEndpoint, error)
	FetchPeerByKeyId(iContext context.Context, iPublicKeyId models.PublicKeyId) (models.Peer, error)
}
