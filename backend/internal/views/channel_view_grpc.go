package views

import (
	"backend/internal/controllers"
	message "backend/internal/grpc"
	"context"

	"google.golang.org/grpc"
)

type ChannelViewGrpc struct {
	message.UnimplementedChannelGrpcServer
	controller *controllers.ChannelController
}

func MakeChannelViewGrpc(
	controller *controllers.ChannelController,
) (*ChannelViewGrpc, error) {
	return &ChannelViewGrpc{
		controller: controller,
	}, nil
}

func (v *ChannelViewGrpc) Register(grpcServer *grpc.Server) {
	message.RegisterChannelGrpcServer(grpcServer, v)
}

func (v *ChannelViewGrpc) GetChannels(context.Context, *message.GetChannelsRequest) (*message.GetChannelsResponse, error) {
	result, err := v.controller.GetChannels()
	response := message.GetChannelsResponse{}
	if err == nil {
		response.Channels = result
	}
	return &response, err
}
