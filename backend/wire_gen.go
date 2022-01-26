// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"backend/internal/controllers"
	"backend/internal/providers"
	"backend/internal/repositories"
	"backend/internal/services"
	"backend/internal/views"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeView() (*views.ConnectionViewGrpc, error) {
	bitcask := providers.MakeBitcask()
	connectionRepositoryBitcask := repositories.MakeConnectionRepositoryBitcask(bitcask)
	string2 := _wireStringValue
	connectionStorageService := services.MakeConnectionStorageService(connectionRepositoryBitcask, string2)
	userRepositoryBitcask := repositories.MakeUserRepositoryBitcask(bitcask)
	userStorageService := services.MakeUserStorageService(userRepositoryBitcask)
	connectionController := controllers.MakeController(connectionStorageService, userStorageService)
	connectionViewGrpc, err := views.MakeConnectionViewGrpc(connectionController)
	if err != nil {
		return nil, err
	}
	return connectionViewGrpc, nil
}

var (
	_wireStringValue = "./"
)

func InitializeUserView() (*views.UserViewGrpc, error) {
	bitcask := providers.MakeBitcask()
	userRepositoryBitcask := repositories.MakeUserRepositoryBitcask(bitcask)
	userStorageService := services.MakeUserStorageService(userRepositoryBitcask)
	wallet := providers.MakeHyperledgerWallet()
	enrollUserHyperledgerSerivce := services.MakeEnrollUserHyperledgerService(wallet)
	userController := controllers.MakeUserController(userStorageService, enrollUserHyperledgerSerivce)
	userViewGrpc, err := views.MakeUserViewGrpc(userController)
	if err != nil {
		return nil, err
	}
	return userViewGrpc, nil
}

func InitializeChannelView() (*views.ChannelViewGrpc, error) {
	channelHyperledgerService := services.MakeChannelHyperledgerService()
	channelController := controllers.MakeChannelController(channelHyperledgerService)
	channelViewGrpc, err := views.MakeChannelViewGrpc(channelController)
	if err != nil {
		return nil, err
	}
	return channelViewGrpc, nil
}

// wire.go:

var repositorySet = wire.NewSet(repositories.MakeConnectionRepositoryBitcask, wire.Bind(
	new(repositories.IConnectionRepository),
	new(*repositories.ConnectionRepositoryBitcask),
),
)

var userRepositorySet = wire.NewSet(repositories.MakeUserRepositoryBitcask, wire.Bind(
	new(repositories.UserRepositoryI),
	new(*repositories.UserRepositoryBitcask),
),
)
