package material_contract_service

import (
	"backend/internal/drivers"
	"encoding/json"
	"reflect"
)

type MaterialVerificationService struct {
	driver drivers.SmartContractDriverI
}

func (s MaterialVerificationService) Verify(
	iMaterial material,
) (bool, error) {
	nodeJson, err := s.driver.Query(
		"GetMaterial",
		iMaterial.Id,
	)
	if err != nil {
		return false, err
	}

	var materialSc material
	err = json.Unmarshal(nodeJson, &materialSc)
	return reflect.DeepEqual(iMaterial, materialSc), nil
}
