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

func (r *queryResolver) Peers(iCtx context.Context) ([]*Peer, error) {
	peers, err := r.PeersController.ListPeers(iCtx)
	if err != nil {
		return []*Peer{}, err
	}

	ret := make([]*Peer, len(peers))
	for i, peer := range peers {
		parsedPeer := ParsePeer(peer)
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
