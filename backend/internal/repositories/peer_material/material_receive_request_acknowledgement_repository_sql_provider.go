package peer_material_repositories

import "go.uber.org/dig"

func ProvideMaterialReceiveRequestAcknowledgementRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeMaterialReceiveRequestAcknowledgementRepositorySql, dig.As(new(MaterialReceiveRequestAcknowledgementRepositoryI)))
}
