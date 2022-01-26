package views

import (
	"backend/internal/controllers"
	message "backend/internal/grpc"
	"context"

	"google.golang.org/grpc"
)

type ConnectionViewGrpc struct {
	message.UnimplementedConnectionGrpcServer
	controller *controllers.ConnectionController
}

func MakeConnectionViewGrpc(
	controller *controllers.ConnectionController,
) (*ConnectionViewGrpc, error) {
	view := ConnectionViewGrpc{
		controller: controller,
	}
	return &view, nil
}

func (v *ConnectionViewGrpc) Register(grpcServer *grpc.Server) {
	message.RegisterConnectionGrpcServer(grpcServer, v)
}

func (v *ConnectionViewGrpc) GetConnectionProfile(context.Context, *message.GetConnectionProfileRequest) (*message.ConnectionProfile, error) {
	profile, err := v.controller.LoadConnection()
	if err == nil {
		ret := message.ConnectionProfile{
			Data: profile,
		}
		return &ret, err
	} else {
		return nil, err
	}
}

func (v *ConnectionViewGrpc) SaveConnectionProfile(context context.Context, request *message.SaveConnectionProfileRequest) (*message.ConnectionProfile, error) {
	profile, err := v.controller.SaveConnection(request.GetPath())
	if err == nil {
		ret := message.ConnectionProfile{
			Data: profile,
		}
		return &ret, err
	} else {
		return &message.ConnectionProfile{}, err
	}
}

func (v *ConnectionViewGrpc) Connect(context context.Context, request *message.ConnectRequest) (*message.ConnectResponse, error) {
	err := v.controller.Connect()
	return &message.ConnectResponse{}, err
}
