package repositories

import "go.uber.org/dig"

func ProvidePeerKeyRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePeerKeyRepositorySql, dig.As(new(PeerKeyRepositoryI)))
}
