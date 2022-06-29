package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"backend/graphql"
	"backend/internal/common"
	"backend/internal/controllers"
	"backend/internal/drivers"
	"backend/internal/repositories"
	peer_material_repositories "backend/internal/repositories/peer_material"
	"backend/internal/services"
	graph_id_service "backend/internal/services/graph_id"
	material_contract_service "backend/internal/services/material_contract"
	"backend/internal/services/node_contract"
	peer_material_services "backend/internal/services/peer_material"

	pb "backend/internal/grpc"

	"github.com/99designs/gqlgen/handler"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

const defaultPort = "9053"

func initialize(container *dig.Container) {
	var err error
	err = common.ProvideConfig(container)
	if err != nil {
		log.Fatal(err)
	}
	err = common.ProvideSqlDb(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvideNodeRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvideMaterialRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = peer_material_repositories.ProvideMaterialReceiveRequestRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = peer_material_services.ProvideReceiveMaterialRequestRepositoryService(container)
	if err != nil {
		log.Fatal(err)
	}
	err = node_contract.ProvideIdHasher(container)
	if err != nil {
		log.Fatal(err)
	}
	err = services.ProvideUserService(container)
	if err != nil {
		log.Fatal(err)
	}
	err = material_contract_service.ProvideMaterialFetchServiceHl(container)
	if err != nil {
		log.Fatal(err)
	}
	err = controllers.ProvidePeerMaterialController(container)
	if err != nil {
		log.Fatal(err)
	}
	err = peer_material_services.ProvidePeerMaterialServerService(container)
	if err != nil {
		log.Fatal(err)
	}
	err = peer_material_repositories.ProvideMaterialReceiveRequestAcknowledgementRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = material_contract_service.ProvideSignatureOptionGenerator(container)
	if err != nil {
		log.Fatal(err)
	}
	err = services.ProvideMaterialRepositoryService(container)
	if err != nil {
		log.Fatal(err)
	}
	err = peer_material_services.ProvidePeerMaterialClientServiceGrpc(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvideUserRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvideUserKeyRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = drivers.ProvideSmartContractDriverHl(container)
	if err != nil {
		log.Fatal(err)
	}
	err = peer_material_services.ProvidePeerMaterialClientService(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvidePeerRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = node_contract.ProvideSignatureService(container)
	if err != nil {
		log.Fatal(err)
	}
	err = graph_id_service.ProvideIdGeneratorServiceUuid(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvideMaterialRepositoryFactory(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvidePeerKeyRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = material_contract_service.ProvideMaterialTransferServiceHl(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvideUserEndpointRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
	err = repositories.ProvidePeerProtocolRepositorySql(container)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		container := dig.New()
		initialize(container)

		defer wg.Done()
		os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
		port := os.Getenv("PORT")
		if port == "" {
			port = defaultPort
		}

		err := container.Invoke(func(peerMaterialController *controllers.PeerMaterialController) {
			http.Handle("/", handler.Playground("GraphQL playground", "/query"))
			http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(
				graphql.Config{
					Resolvers: &graphql.Resolver{
						MaterialContractController: InitializeMaterialContractController(),
						PeersController:            InitializePeerController(),
						PeerMaterialController:     peerMaterialController,
					},
				},
			)))
			log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
			log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		container := dig.New()
		initialize(container)
		defer wg.Done()
		err := container.Invoke(func(config common.Config, server *peer_material_services.PeerMaterialServerServiceGrpc) {
			lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", config.GrpcServerPort))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			grpcServer := grpc.NewServer()
			pb.RegisterMaterialServiceServer(grpcServer, server)
			grpcServer.Serve(lis)
		})
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}
