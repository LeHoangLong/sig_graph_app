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
}

func MakeMaterialContractController(
	iMaterialServiceFactory material_contract_service.MaterialServiceFactory,
	iRepositoryService services.MaterialRepositoryService,
	iUserKeyRepository repositories.UserKeyRepositoryI,
) MaterialContractController {
	return MaterialContractController{
		materialServiceFactory: iMaterialServiceFactory,
		repositoryService:      iRepositoryService,
		userKeyRepository:      iUserKeyRepository,
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
		userKeyPair.PublicKey,
		material,
	)

	if err != nil {
		return material, fmt.Errorf("material saved to ledger but could not save to db: %s", err.Error())
	}

	return material, nil
}

func (c MaterialContractController) GetMaterialById(
	iCtx context.Context,
	iMaterialId string,
) (models.Material, error) {
	userId, err := services.GetCurrentUserFromContext(iCtx)
	if err != nil {
		return models.Material{}, err
	}

	materialFetchSc := c.materialServiceFactory.BuildMaterialFetchService()
	material, err := materialFetchSc.GetMaterialById(
		iMaterialId,
	)

	doesMaterialBelongToUser, err := c.repositoryService.DoesMaterialBelongToUser(userId, material)
	if err != nil {
		return models.Material{}, err
	}

	keys, err := c.userKeyRepository.FetchUserKeyPairByUser(userId)
	var selectedKey models.PublicKey
	for _, key := range keys {
		if key.PublicKey.Value == material.OwnerPublicKey.Value {
			selectedKey = key.PublicKey
			break
		}
	}
	if doesMaterialBelongToUser {
		if selectedKey.Id != nil { /// This should always happen
			material, err = c.repositoryService.AddMaterialToUser(
				iCtx,
				selectedKey,
				material,
			)
		}
	}

	return material, err
}

func (c MaterialContractController) ListMaterialsOfCurrentUser(
	iCtx context.Context,
) ([]models.Material, error) {
	userId, err := services.GetCurrentUserFromContext(iCtx)
	if err != nil {
		return []models.Material{}, err
	}

	keys, err := c.userKeyRepository.FetchUserKeyPairByUser(userId)
	return c.repositoryService.FetchMaterialsOfUser(iCtx, keys, 0, 10000)
}
