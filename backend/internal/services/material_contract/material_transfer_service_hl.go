package material_contract_service

import (
	"backend/internal/drivers"
	"backend/internal/models"
	"backend/internal/services/node_contract"
	"time"
)

type MaterialTransferService struct {
	signatureService node_contract.SignatureService
	hasher           node_contract.IdHasherI
	driver           drivers.SmartContractDriverI
	fetchService     MaterialFetchServiceI
}

func MakeMaterialTransferServiceHl(
	iSignatureService node_contract.SignatureService,
	iHasher node_contract.IdHasherI,
	iDriver drivers.SmartContractDriverI,
	iFetchService MaterialFetchServiceI,
) MaterialTransferService {
	return MaterialTransferService{
		signatureService: iSignatureService,
		hasher:           iHasher,
		driver:           iDriver,
		fetchService:     iFetchService,
	}
}

func (s MaterialTransferService) TransferMaterial(
	iNewNodeId string,
	iSenderSignature string,
	iOwnerKey models.UserKeyPair,
	iMaterial models.Material,
) (models.Material, models.Material, error) {
	transferTime := models.CustomTime(time.Now())
	newMaterialSc := makeMaterialFromModel(iMaterial)

	currentHashedId := string(s.hasher.Hash(iMaterial.NodeId))
	newMaterialSc.PreviousNodeHashedIds[currentHashedId] = true
	newMaterialSc.CreatedTime = transferTime

	newNodeSignature, err := s.signatureService.CreateNodeSignature(
		iOwnerKey.PrivateKey,
		&newMaterialSc,
	)

	if err != nil {
		return models.Material{}, models.Material{}, err
	}

	_, err = s.driver.CreateTransaction(
		"TransferMaterial",
		iMaterial.NodeId,
		iNewNodeId,
		string(iOwnerKey.PublicKey.Value),
		iSenderSignature,
		newNodeSignature,
	)

	if err != nil {
		return models.Material{}, models.Material{}, err
	}

	transferedMaterial := iMaterial

	newHashedId := string(s.hasher.Hash(iNewNodeId))
	iMaterial.NextNodeHashedIds[newHashedId] = true
	iMaterial.Signature = iSenderSignature
	iMaterial.IsFinalized = true

	transferedMaterial.PreviousNodeHashedIds[currentHashedId] = true
	transferedMaterial.Signature = newNodeSignature
	transferedMaterial.IsFinalized = false

	return iMaterial, transferedMaterial, nil
}
