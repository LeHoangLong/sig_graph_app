package repositories

import "backend/internal/models"

type MaterialRepositoryI interface {
	AddMaterial(iMaterial models.Material) (models.Material, error)
	FetchMaterialsByOwner(iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Material, error)
}
