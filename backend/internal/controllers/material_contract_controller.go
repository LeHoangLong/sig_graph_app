package controllers

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services"
	"context"
	"fmt"
)

type MaterialContractController struct {
	contract   services.MaterialContract
	repository repositories.MaterialRepositoryI
}

func MakeMaterialContractController(
	iContract services.MaterialContract,
	iRepositoy repositories.MaterialRepositoryI,
) MaterialContractController {
	return MaterialContractController{
		contract:   iContract,
		repository: iRepositoy,
	}
}

func (c MaterialContractController) CreateMaterial(
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

	err = c.repository.AddMaterial(
		services.GetUserFromContext(iCtx),
		material,
	)

	if err != nil {
		return material, fmt.Errorf("material saved to ledger but could not save to db: %s", err.Error())
	}

	return material, nil
}

func (c MaterialContractController) GetMaterialById(
	iMaterialId string,
) (models.Material, error) {
	material, err := c.contract.GetMaterialById(
		iMaterialId,
	)
	return material, err
}

func (c MaterialContractController) ListMaterials(
	iCtx context.Context,
) ([]models.Material, error) {
	materials, err := c.repository.FetchMaterials(services.GetUserFromContext(iCtx))
	if err != nil {
		return []models.Material{}, err
	}

	return materials, nil
}
