package controllers

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services"
	material_contract_service "backend/internal/services/material_contract"
	"context"
	"fmt"
)

type MaterialContractController struct {
	materialServiceFactory material_contract_service.MaterialServiceFactory
	repositoryService      services.MaterialRepositoryService
	userKeyRepository      repositories.UserKeyRepositoryI
	peerKeyRepository      repositories.PeerKeyRepositoryI
}

func MakeMaterialContractController(
	iMaterialServiceFactory material_contract_service.MaterialServiceFactory,
	iRepositoryService services.MaterialRepositoryService,
	iUserKeyRepository repositories.UserKeyRepositoryI,
	iPeerKeyRepository repositories.PeerKeyRepositoryI,
) MaterialContractController {
	return MaterialContractController{
		materialServiceFactory: iMaterialServiceFactory,
		repositoryService:      iRepositoryService,
		userKeyRepository:      iUserKeyRepository,
		peerKeyRepository:      iPeerKeyRepository,
	}
}

func (c MaterialContractController) CreateMaterialForCurrentUser(
	iCtx context.Context,
	iName string,
	iUnit string,
	iQuantity string,
) (models.Material, error) {
	userId, err := services.GetCurrentUserFromContext(iCtx)
	if err != nil {
		return models.Material{}, err
	}

	userKeyPair, err := c.userKeyRepository.FetchDefaultUserKeyPair(userId)
	if err != nil {
		return models.Material{}, err
	}

	quantity, err := models.NewCustomDecimalFromString(iQuantity)
	if err != nil {
		return models.Material{}, fmt.Errorf("could not parse quantity: %s", err.Error())
	}

	materialCreateSc := c.materialServiceFactory.BuildMaterialCreateService(userKeyPair)
	material, err := materialCreateSc.CreateMaterial(
		iCtx,
		iName,
		iUnit,
		quantity,
		userKeyPair,
	)

	if err != nil {
		return models.Material{}, err
	}

	material, err = c.repositoryService.AddMaterialToUser(
		iCtx,
		userId,
		userKeyPair.PublicKey,
		material,
	)

	if err != nil {
		return material, fmt.Errorf("material saved to ledger but could not save to db: %s", err.Error())
	}

	return material, nil
}

func (c MaterialContractController) GetMaterialByNodeId(
	iCtx context.Context,
	iUserId models.UserId,
	iMaterialNodeId string,
) (models.Material, error) {
	materialFetchSc := c.materialServiceFactory.BuildMaterialFetchService()
	material, err := materialFetchSc.GetMaterialById(
		iMaterialNodeId,
	)

	keys, err := c.userKeyRepository.FetchUserKeyPairByUser(iUserId)
	var selectedKey models.PublicKey
	for _, key := range keys {
		if key.PublicKey.Value == material.OwnerPublicKey.Value {
			selectedKey = key.PublicKey
			break
		}
	}

	if selectedKey.Id == nil {
		peerKeys, err := c.peerKeyRepository.CreateOrFetchPeerKeysByValue(
			iCtx,
			iUserId,
			[]string{material.OwnerPublicKey.Value},
		)
		if err != nil {
			return models.Material{}, err
		}
		selectedKey = peerKeys[0].PublicKey
	}

	material.OwnerPublicKey = selectedKey
	namespace := services.UserIdToNamespace(iUserId)
	materials, err := c.repositoryService.SaveMaterialsIgnoreId(
		iCtx,
		namespace,
		[]models.Material{material},
	)

	if err != nil {
		return models.Material{}, err
	}

	return materials[0], nil
}

func (c MaterialContractController) ListMaterialsOfCurrentUser(
	iCtx context.Context,
) ([]models.Material, error) {
	userId, err := services.GetCurrentUserFromContext(iCtx)
	if err != nil {
		return []models.Material{}, err
	}

	keys, err := c.userKeyRepository.FetchUserKeyPairByUser(userId)
	return c.repositoryService.FetchMaterialsOfUser(iCtx, userId, keys, 0, 10000)
}
