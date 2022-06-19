package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"backend/internal/services"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateMaterial(ctx context.Context, name string, unit string, quantity string) (*Material, error) {
	material, err := r.MaterialContractController.CreateMaterialForCurrentUser(
		ctx,
		name,
		unit,
		quantity,
	)

	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return nil, err
	}

	parsedMaterial := ParseMaterial(material)
	return &parsedMaterial, nil
}

func (r *mutationResolver) TransferMaterial(ctx context.Context, materialID int, peerID int, peerPublicKeyID int, relatedMaterialsID []int) (*ReceiveMaterialRequestResponse, error) {
	request, err := r.PeerMaterialController.SendRequest(
		ctx,
		peerID,
		peerPublicKeyID,
		materialID,
		relatedMaterialsID,
	)
	if err != nil {
		return nil, err
	}

	parsedRequest := ParseReceiveMaterialRequestRequest(request)
	parsedResponse := ReceiveMaterialRequestResponse{
		Request:  &parsedRequest,
		Accepted: true,
	}
	return &parsedResponse, nil
}

func (r *queryResolver) MaterialByNodeID(ctx context.Context, nodeID string) (*Material, error) {
	material, err := r.MaterialContractController.GetMaterialByNodeId(ctx, nodeID)
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

func (r *queryResolver) PendingReceivedTransferMaterialRequests(ctx context.Context) ([]*ReceiveMaterialRequestRequest, error) {
	senderId, err := services.GetCurrentUserFromContext(ctx)
	if err != nil {
		return []*ReceiveMaterialRequestRequest{}, err
	}
	requests, err := r.PeerMaterialController.FetchReceivedPendingMaterialReceiveRequests(
		ctx,
		senderId,
	)
	if err != nil {
		return []*ReceiveMaterialRequestRequest{}, err
	}
	parsedRequests := []*ReceiveMaterialRequestRequest{}
	for i := range requests {
		parsedRequest := ParseReceiveMaterialRequestRequest(requests[i])
		parsedRequests = append(parsedRequests, &parsedRequest)
	}

	return parsedRequests, nil
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
func (r *queryResolver) ReceivedTransferMaterialRequests(ctx context.Context) ([]*ReceiveMaterialRequestRequest, error) {
	userId, err := services.GetCurrentUserFromContext(ctx)
	if err != nil {
		return []*ReceiveMaterialRequestRequest{}, nil
	}
	requests, err := r.PeerMaterialController.FetchReceivedPendingMaterialReceiveRequests(
		ctx,
		userId,
	)

	parsedRequests := []*ReceiveMaterialRequestRequest{}
	for i := range requests {
		parsedRequest := ParseReceiveMaterialRequestRequest(requests[i])
		parsedRequests = append(parsedRequests, &parsedRequest)
	}

	return parsedRequests, nil
}
