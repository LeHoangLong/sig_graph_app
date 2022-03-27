package repositories

import (
	"backend/internal/models"
)

type PeerRepositoryI interface {
	FetchPeers(iUserId int) ([]models.Peer, error)
}
