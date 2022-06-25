package peer_material_repositories

import "go.uber.org/dig"

func ProvideMaterialReceiveRequestRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakeMaterialReceiveRequestRepositorySql, dig.As(new(MaterialReceiveRequestRepositoryI)))
}
