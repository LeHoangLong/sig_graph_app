package graph_id_service

import (
	"context"

	"github.com/google/uuid"
)

type IdGeneratorServiceUuidPrefix string

type IdGeneratorServiceUuid struct {
	prefix IdGeneratorServiceUuidPrefix
}

func MakeIdGeneratorServiceUuid(
	iPrefix IdGeneratorServiceUuidPrefix,
) IdGeneratorServiceUuid {
	return IdGeneratorServiceUuid{
		prefix: iPrefix,
	}
}

func (s IdGeneratorServiceUuid) NewId(
	iContext context.Context,
) string {
	return string(s.prefix) + "/" + uuid.New().String()
}
