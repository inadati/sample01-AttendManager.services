package main

import (
	"log"
	"net/http"
	"os"
	graphql_cron "root"

	"github.com/99designs/gqlgen/handler"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {

	// dotenv file loading
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql_cron.NewExecutableSchema(graphql_cron.Config{Resolvers: &graphql_cron.Resolver{}})))

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Printf("graphql endpoint is http://localhost:%s/query", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
