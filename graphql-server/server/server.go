package main

import (
	gqlgen "root"

	"log"
	"net/http"
	"os"
	"root/middlewares/auth"
	"root/middlewares/responseHeader"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
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

	router := chi.NewRouter()

	router.Use(responseHeader.Middleware())
	router.Use(auth.Middleware())

	// router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &gqlgen.Resolver{}})))

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Printf("graphql endpoint is http://localhost:%s/query", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
