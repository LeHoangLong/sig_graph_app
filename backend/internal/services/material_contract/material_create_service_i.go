package material_contract_service

import (
	"backend/internal/models"
	"context"
)

type MaterialCreateServiceI interface {
	/// the returned material will have nil id (not yet saved into db)
	CreateMaterial(
		iContext context.Context,
		iName string,
		iUnit string,
		iQuantity models.CustomDecimal,
		iOwnerKey models.UserKeyPair,
	) (models.Material, error)
}
