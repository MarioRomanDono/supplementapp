package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	service := createSupplementService()
	server := NewServer(service)
	log.Fatal(http.ListenAndServe(":8080", server))
}

func createSupplementService() *supplement.SupplementService {
	repo := createMongoDBSupplementRepository()
	service := supplement.NewSupplementService(repo)

	return service
}

func createMongoDBSupplementRepository() *supplement.MongoDBSupplementRepository {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatalf("could not connect to mongodb: %v", err)
	}

	collection := client.Database(os.Getenv("MONGODB_DB")).Collection("supplements")
	repo := supplement.NewMongoDBSupplementRepository(collection)
	return repo
}

func NewServer(service *supplement.SupplementService) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, service)
	return mux
}
