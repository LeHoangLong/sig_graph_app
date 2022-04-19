package graph_id_service

import "context"

type IdGeneratorServiceI interface {
	NewId(iContext context.Context) string
}
