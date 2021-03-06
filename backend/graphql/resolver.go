package graphql

import "backend/internal/controllers"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MaterialContractController controllers.MaterialContractController
	PeersController            controllers.PeersController
	PeerMaterialController     *controllers.PeerMaterialController
}
