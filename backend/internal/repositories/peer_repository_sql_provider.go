package repositories

import "go.uber.org/dig"

func ProvidePeerRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(MakePeerRepositorySql, dig.As(new(PeerRepositoryI)))
}
