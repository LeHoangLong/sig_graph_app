package main

import (
	"log"
	"net/http"
	"os"

	"backend/graphql"

	"github.com/99designs/gqlgen/handler"
)

const defaultPort = "9053"

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(
		graphql.Config{
			Resolvers: &graphql.Resolver{
				MaterialContractController: InitializeMaterialContractController(),
			},
		},
	)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
