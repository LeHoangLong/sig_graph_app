package repositories

import (
	"backend/internal/models"
	"context"
)

type NodeRepositoryI interface {
	CreateNode(iNode models.Node) (models.Node, error)
	FetchNodesByOwnerKey(iContext context.Context, iNamespace string, iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Node, error)
	FetchNodesById(iId []models.NodeId) ([]models.Node, error)
	FetchNodesByNodeId(
		iContext context.Context,
		iNamespace string,
		iId map[string]bool,
	) (map[models.NodeId]models.Node, error)
	UpsertNodesById(
		iContext context.Context,
		iNodes map[models.NodeId]models.Node,
	) (map[models.NodeId]models.Node, error)
	UpsertNodesByNodeIdAndNamespace(
		iContext context.Context,
		iNamespace string,
		iNodes map[string]models.Node,
	) (map[string]models.Node, error)
}
