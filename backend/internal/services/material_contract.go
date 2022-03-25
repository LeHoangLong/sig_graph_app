package services

import (
	"backend/internal/drivers"
	"backend/internal/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MaterialContract struct {
	driver           drivers.SmartContractDriverI
	signatureService GraphContractSignature
	publicKey        PublicKey
}

type materialSC struct {
	Name        string               `json:"Name"`
	Unit        string               `json:"Unit"`
	Quantity    models.CustomDecimal `json:"Quantity"`
	CreatedTime models.CustomTime    `json:"CreatedTime"`
}

func MakeMaterialContract(
	iDriver drivers.SmartContractDriverI,
	iSignatureService GraphContractSignature,
	iPublicKey PublicKey,
) MaterialContract {
	return MaterialContract{
		iDriver,
		iSignatureService,
		iPublicKey,
	}
}

func (c MaterialContract) CreateMaterial(
	iName string,
	iUnit string,
	iQuantity models.CustomDecimal,
) (models.Material, error) {

	material := materialSC{
		Name:        iName,
		Unit:        iUnit,
		Quantity:    iQuantity,
		CreatedTime: models.CustomTime(time.Now()),
	}

	nodeId := uuid.New().String()
	signature, err := c.signatureService.CreateNodeSignature(
		nodeId,
		false,
		material,
	)

	if err != nil {
		return models.Material{}, err
	}

	createdTime := models.CustomTime(time.Now())
	_, err = c.driver.CreateTransaction(
		"CreateMaterial",
		nodeId,
		iName,
		iUnit,
		iQuantity.String(),
		string(c.publicKey),
		createdTime.String(),
		signature,
	)

	if err != nil {
		return models.Material{}, err
	}

	newMaterial := models.NewMaterial(
		nodeId,
		iName,
		iQuantity,
		iUnit,
		createdTime,
	)
	return newMaterial, err
}

func (c MaterialContract) GetMaterialById(
	iNodeId string,
) (models.Material, error) {
	materialJson, err := c.driver.Query(
		"GetMaterial",
		iNodeId,
	)
	if err != nil {
		return models.Material{}, err
	}

	var materialSc materialSC
	err = json.Unmarshal(materialJson, &materialSc)
	if err != nil {
		return models.Material{}, err
	}

	return models.NewMaterial(
		iNodeId,
		materialSc.Name,
		materialSc.Quantity,
		materialSc.Unit,
		materialSc.CreatedTime,
	), nil
}
