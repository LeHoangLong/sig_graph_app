package peer_material_repositories

import "go.uber.org/dig"

func ProvidePendingMaterialReceiveRequestResponseRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePendingMaterialReceiveRequestResponseRepositorySql, dig.As(new(PendingMaterialReceiveRequestResponseRepositoryI)))
}
