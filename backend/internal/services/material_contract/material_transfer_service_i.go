package material_contract_service

import "backend/internal/models"

type MaterialTransferServiceI interface {
	/// first one is the old material with new node id added
	/// second one is the new material
	TransferMaterial(
		iNewNodeId string,
		iSenderSignature string,
		iOwnerKey models.UserKeyPair,
		iMaterial models.Material,
	) (models.Material, models.Material, error)
}
