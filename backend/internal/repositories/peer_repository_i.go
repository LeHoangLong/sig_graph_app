package repositories

import (
	"backend/internal/models"
	"context"
)

type PeerRepositoryI interface {
	FetchPeers(iUserId int) ([]models.Peer, error)
	FetchPeerEndPoints(iContext context.Context, iPeerId int) ([]models.PeerEndpoint, error)
}
