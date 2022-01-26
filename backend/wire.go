// +build wireinject

package main

import (
	"backend/internal/controllers"
	"backend/internal/providers"
	"backend/internal/repositories"
	"backend/internal/services"
	"backend/internal/views"

	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	repositories.MakeConnectionRepositoryBitcask,
	wire.Bind(
		new(repositories.IConnectionRepository),
		new(*repositories.ConnectionRepositoryBitcask),
	),
)

var userRepositorySet = wire.NewSet(
	repositories.MakeUserRepositoryBitcask,
	wire.Bind(
		new(repositories.UserRepositoryI),
		new(*repositories.UserRepositoryBitcask),
	),
)

func InitializeView() (*views.ConnectionViewGrpc, error) {
	wire.Build(
		providers.MakeBitcask,
		views.MakeConnectionViewGrpc,
		controllers.MakeController,
		services.MakeConnectionStorageService,
		services.MakeUserStorageService,
		repositorySet,
		userRepositorySet,
		wire.Value("./assets/connection.json"),
	)
	return &views.ConnectionViewGrpc{}, nil
}

func InitializeUserView() (*views.UserViewGrpc, error) {
	wire.Build(
		providers.MakeBitcask,
		providers.MakeHyperledgerWallet,
		views.MakeUserViewGrpc,
		controllers.MakeUserController,
		services.MakeEnrollUserHyperledgerService,
		services.MakeUserStorageService,
		userRepositorySet,
	)
	return &views.UserViewGrpc{}, nil
}

func InitializeChannelView() (*views.ChannelViewGrpc, error) {
	wire.Build(
		views.MakeChannelViewGrpc,
		controllers.MakeChannelController,
		services.MakeChannelHyperledgerService,
	)
	return &views.ChannelViewGrpc{}, nil
}
