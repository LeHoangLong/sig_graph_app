package controllers

import (
	peer_material_services "backend/internal/services/peer_material"

	"go.uber.org/dig"
)

func ProvidePeerMaterialController(
	iContainer *dig.Container,
) error {
	err := iContainer.Provide(MakePeerMaterialController)
	if err != nil {
		return err
	}
	err = iContainer.Provide(func(iController *PeerMaterialController) peer_material_services.ReceiveMaterialRequestReceivedHandler {
		return iController
	})
	if err != nil {
		return err
	}
	err = iContainer.Provide(func(iController *PeerMaterialController) peer_material_services.ReceiveMaterialResponseHandler {
		return iController
	})
	if err != nil {
		return err
	}
	return nil
}
