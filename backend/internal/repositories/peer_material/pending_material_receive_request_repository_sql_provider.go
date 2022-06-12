package peer_material_repositories

import "go.uber.org/dig"

func ProvidePendingMaterialReceiveRequestRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePendingMaterialReceiveRequestRepositorySql, dig.As(new(PendingMaterialReceiveRequestRepositoryI)))
}
