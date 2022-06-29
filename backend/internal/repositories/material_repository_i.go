package repositories

import (
	"backend/internal/models"
	"context"
)

type MaterialRepositoryI interface {
	/// returns AlreadyExists error if duplicate found
	AddMaterial(iMaterial models.Material) (models.Material, error)
	/// upsert based on material.Id. Error if conflict on other columns
	/// Guarantees output is identical to input if no error
	UpsertMaterialsByIds(iContext context.Context, iMaterials map[models.NodeId]models.Material) (map[models.NodeId]models.Material, error)
	/// upsert based on material.NodeId and material.Namespace. Error if conflict on other columns
	/// Guarantees output size == input size and identical to input (except namespace)
	/// iMaterials namespace is ignored and output will have namespace == iNamespace
	UpsertMaterialsByNodeIdsAndNamespace(iContext context.Context, iNamespace string, iMaterials map[string]models.Material) (map[string]models.Material, error)

	FetchMaterialsByOwner(iContext context.Context, iNamespace string, iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Material, error)
	FetchMaterialById(iContext context.Context, iId models.NodeId) (models.Material, error)
	/// ignore ids that don't exist
	FetchMaterialsByNodeId(iContext context.Context, iNamespace string, iMaterialNodeId map[string]bool) (map[models.NodeId]models.Material, error)
	/// return NotFound if any id in iIds is not found
	FetchMaterialsById(iContext context.Context, iIds map[models.NodeId]bool) (map[models.NodeId]models.Material, error)
}
