package controllers

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services"
	"context"
)

type PeersController struct {
	r repositories.PeerRepositoryI
}

func MakePeersController(
	iRepository repositories.PeerRepositoryI,
) PeersController {
	return PeersController{
		r: iRepository,
	}
}

func (c PeersController) ListPeers(
	iContext context.Context,
) ([]models.Peer, error) {
	userId, err := services.GetCurrentUserFromContext(iContext)
	if err != nil {
		return []models.Peer{}, err
	}
	return c.r.FetchPeers(userId)
}
