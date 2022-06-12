package repositories

import (
	"backend/internal/models"
	"context"
)

type MaterialRepositoryI interface {
	AddMaterial(iMaterial models.Material) (models.Material, error)
	FetchMaterialsByOwner(iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Material, error)
	FetchMaterialById(iContext context.Context, iId int) (models.Material, error)
}
