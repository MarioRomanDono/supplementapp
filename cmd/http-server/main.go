package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/marioromandono/supplementapp/internal/supplement"
	"github.com/marioromandono/supplementapp/internal/supplement/persistence/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// TODO: This code is not ready for production as it misses graceful shutdown, closing database connections...
	service := createSupplementService()
	server := NewServer(service)
	log.Fatal(http.ListenAndServe(":8080", server))
}

func createSupplementService() *supplement.SupplementService {
	repo := createPostgresSupplementRepository()
	service := supplement.NewSupplementService(repo)

	return service
}

func createPostgresSupplementRepository() *postgres.PostgresSupplementRepository {
	db, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("could not create postgres pool: %v", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	return postgres.NewSupplementRepository(db)
}

func NewServer(service *supplement.SupplementService) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, service)
	return mux
}
