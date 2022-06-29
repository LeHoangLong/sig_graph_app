package services

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services/node_contract"
	"container/list"
	"context"
	"fmt"
	"strconv"
)

type MaterialRepositoryService struct {
	repositoryFactory  repositories.MaterialRepositoryFactory
	userKeyRepository  repositories.UserKeyRepositoryI
	materialRepository repositories.MaterialRepositoryI
	idHasher           node_contract.IdHasherI
}

func MakeMaterialRepositoryService(
	iRepositoryFactory repositories.MaterialRepositoryFactory,
	iUserKeyRepository repositories.UserKeyRepositoryI,
	iMaterialRepository repositories.MaterialRepositoryI,
	iIdHasher node_contract.IdHasherI,
) MaterialRepositoryService {
	return MaterialRepositoryService{
		repositoryFactory:  iRepositoryFactory,
		userKeyRepository:  iUserKeyRepository,
		materialRepository: iMaterialRepository,
		idHasher:           iIdHasher,
	}
}

func UserIdToNamespace(
	iUserId models.UserId,
) string {
	return strconv.Itoa(int(iUserId))
}

func (s MaterialRepositoryService) AddMaterialToUser(
	iContext context.Context,
	iUserId models.UserId,
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
		namespace := UserIdToNamespace(iUserId)
		iMaterial.Node.Namespace = &namespace
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
	iUserId models.UserId,
	iUserKeys []models.UserKeyPair,
	iMinId int,
	iLimit int,
) ([]models.Material, error) {
	var err error
	ret := []models.Material{}
	namespace := UserIdToNamespace(iUserId)
	err = s.repositoryFactory.GetRepository(iContext, func(iRepository repositories.MaterialRepositoryI) error {
		for _, key := range iUserKeys {
			materials, err := iRepository.FetchMaterialsByOwner(
				iContext,
				namespace,
				key.PublicKey,
				iMinId,
				iLimit,
			)
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
	iUserId models.UserId,
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

/// save materials into namespace
/// update if material with same NodeId and Namespace already exists
/// material.Id is ignored
/// Guarantees that output size == input size
/// All output materials will be in iNamespace
/// output materials will have their id filled
func (s *MaterialRepositoryService) SaveMaterialsIgnoreId(
	iContext context.Context,
	iNamespace string,
	iMaterials []models.Material,
) ([]models.Material, error) {
	materialsByNodeId := map[string]models.Material{}
	for _, material := range iMaterials {
		material.Id = nil
		material.Namespace = &iNamespace
		materialsByNodeId[material.NodeId] = material
	}
	savedMaterials, err := s.materialRepository.UpsertMaterialsByNodeIdsAndNamespace(
		iContext,
		iNamespace,
		materialsByNodeId,
	)

	if err != nil {
		return []models.Material{}, err
	}

	ret := []models.Material{}
	for nodeId := range savedMaterials {
		ret = append(ret, savedMaterials[nodeId])
	}
	return ret, nil
}

/// all materials must have namespace iNamespace. Else returns InvalidArguments error
func (s *MaterialRepositoryService) LinkMaterials(
	iContext context.Context,
	iNamespace string,
	iMaterials map[models.NodeId]models.Material,
) (map[models.NodeId]models.Material, error) {
	nodeIds := map[string]bool{}
	for id := range iMaterials {
		if iMaterials[id].Namespace == nil || *iMaterials[id].Namespace != iNamespace {
			return map[models.NodeId]models.Material{}, common.InvalidArgument
		}
		nodeIds[iMaterials[id].NodeId] = true
	}

	savedMaterials, err := s.materialRepository.FetchMaterialsByNodeId(
		iContext,
		iNamespace,
		nodeIds,
	)

	if err != nil {
		return map[models.NodeId]models.Material{}, err
	}

	materialsByHashedNodeId := map[string]models.Material{}
	for id := range savedMaterials {
		hashedId := s.idHasher.Hash(savedMaterials[id].NodeId)
		materialsByHashedNodeId[hashedId] = savedMaterials[id]
	}

	linkedMaterials := map[models.NodeId]models.Material{}
	for id := range iMaterials {
		material := iMaterials[id]
		for parentHashedId := range material.PreviousNodeHashedIds {
			if node, ok := materialsByHashedNodeId[parentHashedId]; ok {
				material.ParentIds[*node.Id] = true
			}
		}
		for childHashedId := range material.NextNodeHashedIds {
			if node, ok := materialsByHashedNodeId[childHashedId]; ok {
				material.ChildrenIds[*node.Id] = true
			}
		}
		linkedMaterials[*material.Id] = material
	}

	ret, err := s.materialRepository.UpsertMaterialsByIds(
		iContext,
		linkedMaterials,
	)
	if err != nil {
		return map[models.NodeId]models.Material{}, err
	}

	return ret, nil
}

func (s *MaterialRepositoryService) FetchMaterialById(
	iContext context.Context,
	iMaterialId models.NodeId,
) (models.Material, error) {
	return s.materialRepository.FetchMaterialById(iContext, iMaterialId)
}

func (s *MaterialRepositoryService) FetchMaterialByNodeId(
	iContext context.Context,
	iNamespace string,
	iMaterialNodeId string,
) (models.Material, error) {
	materials, err := s.materialRepository.FetchMaterialsByNodeId(
		iContext,
		iNamespace,
		map[string]bool{iMaterialNodeId: true},
	)
	if err != nil {
		return models.Material{}, err
	}

	if len(materials) == 0 {
		return models.Material{}, common.NotFound
	}

	material := models.Material{}
	for id := range materials {
		material = materials[id]
		break
	}

	return material, nil
}

func (s *MaterialRepositoryService) FetchMaterialsAndRelated(
	iContext context.Context,
	iMainMaterialId models.NodeId,
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
		material, err := s.materialRepository.FetchMaterialById(iContext, id.Value.(models.NodeId))
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
