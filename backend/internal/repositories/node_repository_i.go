package repositories

import (
	"backend/internal/models"
)

type NodeRepositoryI interface {
	CreateNode(iNode models.Node) (models.Node, error)
	FetchNodesByOwnerKey(iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Node, error)
}
