package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"fmt"
)

type MaterialRepositoryService struct {
	repository        repositories.MaterialRepositoryI
	userKeyRepository repositories.UserKeyRepositoryI
}

func MakeMaterialRepositoryService(
	iRepository repositories.MaterialRepositoryI,
	iUserKeyRepository repositories.UserKeyRepositoryI,
) MaterialRepositoryService {
	return MaterialRepositoryService{
		repository:        iRepository,
		userKeyRepository: iUserKeyRepository,
	}
}

func (s MaterialRepositoryService) AddOrUpdateMaterialToUser(
	iUserId int,
	iMaterial models.Material,
) (models.Material, error) {
	keys, err := s.userKeyRepository.FetchUserKeyPairByUser(iUserId)
	if err != nil {
		return models.Material{}, err
	}

	if len(keys) == 0 {
		return models.Material{}, fmt.Errorf("no public key created yet")
	}

	defaultKey := keys[0]
	for _, key := range keys {
		if key.IsDefault {
			defaultKey = key
			break
		}
	}

	err = s.repository.AddOrUpdateMaterial(
		defaultKey.PublicKey.Id,
		iMaterial,
	)

	if err != nil {
		return models.Material{}, err
	}

	return iMaterial, nil
}

func (s MaterialRepositoryService) FetchMaterialsOfUser(
	iUserId int,
) ([]models.Material, error) {

	keys, err := s.userKeyRepository.FetchUserKeyPairByUser(iUserId)
	if err != nil {
		return []models.Material{}, err
	}

	ret := []models.Material{}
	for _, key := range keys {
		materials, err := s.repository.FetchMaterials(key.PublicKey.Id)
		if err != nil {
			return []models.Material{}, err
		}

		ret = append(ret, materials...)
	}

	return ret, nil
}

func (s MaterialRepositoryService) DoesMaterialBelongToUser(
	iUserId int,
	iMaterial models.Material,
) (bool, error) {

	keys, err := s.userKeyRepository.FetchUserKeyPairByUser(iUserId)
	if err != nil {
		return false, err
	}

	for _, key := range keys {
		if key.PublicKey.Value == iMaterial.OwnerPublicKey {
			return true, nil
		}
	}

	return false, nil
}
