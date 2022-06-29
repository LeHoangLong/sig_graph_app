package repositories

import "go.uber.org/dig"

func ProvidePeerProtocolRepositorySql(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(
		MakePeerProtocolRepositorySql,
		dig.As(new(PeerProtocolRepositoryI)),
	)
}
