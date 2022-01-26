package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func main() {
	view, err := InitializeView()
	if err != nil {
		panic("Could not initialize connection view")
	}
	userView, err := InitializeUserView()
	if err != nil {
		panic("Could not initialize user view")
	}

	channelView, err := InitializeChannelView()
	if err != nil {
		panic("Could not initialize channel view")
	}

	grpcServer := grpc.NewServer()

	view.Register(grpcServer)
	userView.Register(grpcServer)
	channelView.Register(grpcServer)

	lis, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Printf("Could not listen to localhost:8000: %s. Exiting", err.Error())
		return
	}

	grpcServer.Serve(lis)
}
