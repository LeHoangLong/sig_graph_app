package material_contract_service

import (
	"backend/internal/models"
	graph_id_service "backend/internal/services/graph_id"
	"backend/internal/services/node_contract"
	"context"
)

type SignatureOptionGenerator struct {
	signatureService node_contract.SignatureService
	idGenerator      graph_id_service.IdGeneratorServiceI
	hasher           node_contract.IdHasherI
}

func MakeSignatureOptionGenerator(
	iSignatureService node_contract.SignatureService,
	iIdGenerator graph_id_service.IdGeneratorServiceI,
	IHasher node_contract.IdHasherI,
) SignatureOptionGenerator {
	return SignatureOptionGenerator{
		signatureService: iSignatureService,
		idGenerator:      iIdGenerator,
		hasher:           IHasher,
	}
}

func (s SignatureOptionGenerator) GenerateSignatureOptionsForMaterialTransfer(
	iContext context.Context,
	iNumberOfOptions int,
	iModel models.Material,
	iOwnerKey models.UserKeyPair,
) ([]models.SignatureOption, error) {
	ret := make([]models.SignatureOption, 0, iNumberOfOptions)
	for i := 0; i < iNumberOfOptions; i++ {
		materialSc := makeMaterialFromModel(iModel)
		var id string
		for {
			id = s.idGenerator.NewId(iContext)
			hashedId := s.hasher.Hash(id)
			if _, ok := materialSc.NextNodeHashedIds[hashedId]; !ok {
				materialSc.NextNodeHashedIds[hashedId] = true
				break
			}
		}

		signature, err := s.signatureService.CreateNodeSignature(
			iOwnerKey.PrivateKey,
			&materialSc,
		)
		if err != nil {
			return []models.SignatureOption{}, err
		}

		option := models.MakeSignatureOption(id, signature)
		ret = append(ret, option)
	}
	return ret, nil
}
