package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"fmt"
)

type MaterialRepositoryService struct {
	repositoryFactory repositories.MaterialRepositoryFactory
	userKeyRepository repositories.UserKeyRepositoryI
}

func MakeMaterialRepositoryService(
	iRepositoryFactory repositories.MaterialRepositoryFactory,
	iUserKeyRepository repositories.UserKeyRepositoryI,
) MaterialRepositoryService {
	return MaterialRepositoryService{
		repositoryFactory: iRepositoryFactory,
		userKeyRepository: iUserKeyRepository,
	}
}

func (s MaterialRepositoryService) AddMaterialToUser(
	iContext context.Context,
	iUserKey models.PublicKey,
	iMaterial models.Material,
) (models.Material, error) {
	if iUserKey.Id == nil {
		return models.Material{}, fmt.Errorf("iUserKey is not yet saved to db")
	}
	iMaterial.OwnerPublicKey = iUserKey
	var material models.Material
	var err error
	err = s.repositoryFactory.GetRepository(iContext, func(iRepository repositories.MaterialRepositoryI) error {
		material, err = iRepository.AddMaterial(
			iMaterial,
		)
		return err
	})

	if err != nil {
		return models.Material{}, err
	}

	return material, nil
}

func (s MaterialRepositoryService) FetchMaterialsOfUser(
	iContext context.Context,
	iUserKeys []models.UserKeyPair,
	iMinId int,
	iLimit int,
) ([]models.Material, error) {
	var err error
	ret := []models.Material{}
	err = s.repositoryFactory.GetRepository(iContext, func(iRepository repositories.MaterialRepositoryI) error {
		for _, key := range iUserKeys {
			materials, err := iRepository.FetchMaterialsByOwner(key.PublicKey, iMinId, iLimit)
			if err != nil {
				return err
			}

			ret = append(ret, materials...)
		}
		return nil
	})

	if err != nil {
		return []models.Material{}, err
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
		if key.PublicKey.Value == iMaterial.OwnerPublicKey.Value {
			return true, nil
		}
	}

	return false, nil
}
