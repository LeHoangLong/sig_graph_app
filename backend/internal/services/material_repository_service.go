package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"container/list"
	"context"
	"fmt"
)

type MaterialRepositoryService struct {
	repositoryFactory  repositories.MaterialRepositoryFactory
	userKeyRepository  repositories.UserKeyRepositoryI
	materialRepository repositories.MaterialRepositoryI
}

func MakeMaterialRepositoryService(
	iRepositoryFactory repositories.MaterialRepositoryFactory,
	iUserKeyRepository repositories.UserKeyRepositoryI,
	iMaterialRepository repositories.MaterialRepositoryI,
) MaterialRepositoryService {
	return MaterialRepositoryService{
		repositoryFactory:  iRepositoryFactory,
		userKeyRepository:  iUserKeyRepository,
		materialRepository: iMaterialRepository,
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

func (s *MaterialRepositoryService) SaveMaterials(
	iContext context.Context,
	iMaterials []models.Material,
) ([]models.Material, error) {
	savedMaterials := []models.Material{}
	for _, material := range iMaterials {
		savedMaterial, err := s.materialRepository.AddMaterial(material)
		if err != nil {
			return []models.Material{}, err
		}
		savedMaterials = append(savedMaterials, savedMaterial)
	}
	return savedMaterials, nil
}

func (s *MaterialRepositoryService) FetchMaterialById(
	iContext context.Context,
	iMaterialId int,
) (models.Material, error) {
	return s.materialRepository.FetchMaterialById(iContext, iMaterialId)
}

func (s *MaterialRepositoryService) FetchMaterialsAndRelated(
	iContext context.Context,
	iMainMaterialId int,
) (models.Material, []models.Material, error) {
	mainMaterial, err := s.materialRepository.FetchMaterialById(iContext, iMainMaterialId)
	if err != nil {
		return models.Material{}, []models.Material{}, nil
	}

	list := list.New()
	for id := range mainMaterial.ChildrenIds {
		list.PushBack(id)
	}
	for id := range mainMaterial.ParentIds {
		list.PushBack(id)
	}
	relatedMaterial := []models.Material{}
	for list.Len() > 0 {
		id := list.Front()
		material, err := s.materialRepository.FetchMaterialById(iContext, id.Value.(int))
		if err != nil {
			return models.Material{}, []models.Material{}, nil
		}
		relatedMaterial = append(relatedMaterial, material)
		for id := range material.ChildrenIds {
			list.PushBack(id)
		}
		for id := range material.ParentIds {
			list.PushBack(id)
		}
	}
	return mainMaterial, relatedMaterial, nil
}
