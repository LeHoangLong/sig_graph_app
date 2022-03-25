package repositories

import "backend/internal/models"

type MaterialRepositoryI interface {
	AddMaterial(iUserId int, iMaterial models.Material) error
	FetchMaterials(iUserId int) ([]models.Material, error)
}
