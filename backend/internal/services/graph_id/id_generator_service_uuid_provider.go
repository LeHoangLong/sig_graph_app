package graph_id_service

import (
	"backend/internal/common"

	"go.uber.org/dig"
)

func ProvideIdGeneratorServiceUuid(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(func(iConfig common.Config) IdGeneratorServiceUuid {
		return MakeIdGeneratorServiceUuid(IdGeneratorServiceUuidPrefix(iConfig.GraphIdPrefix))
	}, dig.As(new(IdGeneratorServiceI)))
}
