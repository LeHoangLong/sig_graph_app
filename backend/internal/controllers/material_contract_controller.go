package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"context"
	"fmt"
)

type MaterialContractController struct {
	contract          services.MaterialContract
	repositoryService services.MaterialRepositoryService
}

func MakeMaterialContractController(
	iContract services.MaterialContract,
	iRepositoryService services.MaterialRepositoryService,
) MaterialContractController {
	return MaterialContractController{
		contract:          iContract,
		repositoryService: iRepositoryService,
	}
}

func (c MaterialContractController) CreateMaterialForCurrentUser(
	iCtx context.Context,
	iName string,
	iUnit string,
	iQuantity string,
) (models.Material, error) {
	quantity, err := models.NewCustomDecimalFromString(iQuantity)
	if err != nil {
		return models.Material{}, fmt.Errorf("could not parse quantity: %s", err.Error())
	}

	material, err := c.contract.CreateMaterial(
		iName,
		iUnit,
		quantity,
	)

	if err != nil {
		return models.Material{}, err
	}

	userId, err := services.GetCurrentUserFromContext(iCtx)
	if err != nil {
		return models.Material{}, err
	}
	material, err = c.repositoryService.AddOrUpdateMaterialToUser(
		userId,
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

	material, err := c.contract.GetMaterialById(
		iMaterialId,
	)

	doesMaterialBelongToUser, err := c.repositoryService.DoesMaterialBelongToUser(userId, material)
	if err != nil {
		return models.Material{}, err
	}

	if doesMaterialBelongToUser {
		material, err = c.repositoryService.AddOrUpdateMaterialToUser(userId, material)
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
	return c.repositoryService.FetchMaterialsOfUser(userId)
}
