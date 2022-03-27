package repositories

import "backend/internal/models"

type MaterialRepositoryI interface {
	AddOrUpdateMaterial(iPublicKeyId int, iMaterial models.Material) error
	FetchMaterials(iPublicKeyId int) ([]models.Material, error)
}
