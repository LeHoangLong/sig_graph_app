package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CreateMaterial(ctx context.Context, iName string, iUnit string, iQuantity string) (*Material, error) {
	material, err := r.MaterialContractController.CreateMaterialForCurrentUser(
		ctx,
		iName,
		iUnit,
		iQuantity,
	)

	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return nil, err
	}

	parsedMaterial := ParseMaterial(material)
	return &parsedMaterial, nil
}

func (r *queryResolver) Material(ctx context.Context, id string) (*Material, error) {
	material, err := r.MaterialContractController.GetMaterialById(ctx, id)
	if err != nil {
		return nil, err
	}

	parsedMaterial := ParseMaterial(material)
	return &parsedMaterial, nil
}

func (r *queryResolver) Materials(ctx context.Context) ([]*Material, error) {
	materials, err := r.MaterialContractController.ListMaterialsOfCurrentUser(ctx)
	if err != nil {
		return []*Material{}, err
	}

	ret := []*Material{}
	for _, material := range materials {
		parsedMaterial := ParseMaterial(material)
		ret = append(ret, &parsedMaterial)
	}

	return ret, nil
}

func (r *queryResolver) Peers(ctx context.Context) ([]*Peer, error) {
	peers, err := r.PeersController.ListPeers(ctx)
	if err != nil {
		return []*Peer{}, err
	}

	ret := make([]*Peer, len(peers))
	for i, peer := range peers {
		parsedPeer, err := ParsePeer(peer)
		if err != nil {
			return []*Peer{}, err
		}
		ret[i] = &parsedPeer
	}
	return ret, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
