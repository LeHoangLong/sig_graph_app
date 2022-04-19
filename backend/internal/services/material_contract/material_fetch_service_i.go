package material_contract_service

import "backend/internal/models"

type MaterialFetchServiceI interface {
	/// the returned material will have id of -1
	GetMaterialById(
		iNodeId string,
	) (models.Material, error)
}
