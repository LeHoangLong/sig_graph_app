package views

import (
	"backend/internal/controllers"
	message "backend/internal/grpc"
	"context"

	"google.golang.org/grpc"
)

type UserViewGrpc struct {
	message.UnimplementedUserGrpcServer
	controller *controllers.UserController
}

func MakeUserViewGrpc(
	controller *controllers.UserController,
) (*UserViewGrpc, error) {
	view := UserViewGrpc{
		controller: controller,
	}
	return &view, nil
}

func (v *UserViewGrpc) Register(grpcServer *grpc.Server) {
	message.RegisterUserGrpcServer(grpcServer, v)
}

func (v *UserViewGrpc) CreateUser(context context.Context, iMessage *message.CreateUserRequest) (*message.CreateUserResponse, error) {
	err := v.controller.CreateUser(
		iMessage.Username,
		iMessage.Password,
		iMessage.OrganizationMspId,
		iMessage.CertPath,
		iMessage.CertKey,
	)
	return &message.CreateUserResponse{}, err
}

func (v *UserViewGrpc) Login(
	context context.Context,
	iMessage *message.LoginRequest,
) (*message.LoginResponse, error) {
	err := v.controller.Login(
		iMessage.Username,
		iMessage.Password,
	)
	return &message.LoginResponse{}, err
}
