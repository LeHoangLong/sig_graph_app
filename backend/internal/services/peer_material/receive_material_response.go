package peer_material_services

type ReceiveMaterialResponse struct {
	ResponseId        string
	IsRequestAccepted bool
	Message           string
	NewNodeId         string
}
