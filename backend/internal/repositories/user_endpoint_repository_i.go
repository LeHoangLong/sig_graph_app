package repositories

import (
	"backend/internal/models"
	"context"
)

type UserEndpointRepositoryI interface {
	FetchUserEndpointByUserId(
		iContext context.Context,
		iUserId models.UserId,
	) ([]models.UserEndpoint, error)
}
