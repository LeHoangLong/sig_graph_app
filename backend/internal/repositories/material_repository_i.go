package repositories

import (
	"backend/internal/models"
	"context"
)

type MaterialRepositoryI interface {
	AddMaterial(iMaterial models.Material) (models.Material, error)
	FetchMaterialsByOwner(iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Material, error)
	FetchMaterialById(iContext context.Context, iId int) (models.Material, error)
	/// return NotFound if any id in iIds is not found
	FetchMaterialsById(iContext context.Context, iIds map[int]bool) (map[int]models.Material, error)
}
