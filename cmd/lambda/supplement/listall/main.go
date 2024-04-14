package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/marioromandono/supplementapp/internal/supplement"
	"github.com/marioromandono/supplementapp/internal/supplement/persistence/postgres"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LambdaHandler struct {
	service *supplement.SupplementService
}

func main() {
	handler := NewLambdaHandler(createSupplementService())
	lambda.Start(handler.Handle)
}

func (ls *LambdaHandler) Handle(ctx context.Context, r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("REQUEST: ListAll")

	ss, err := ls.service.ListAll(ctx)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	ssJson, err := json.Marshal(ss)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	response := events.APIGatewayV2HTTPResponse{Body: string(ssJson), StatusCode: 200}

	log.Printf("RESPONSE: %v", response)

	return response, nil
}

func NewLambdaHandler(service *supplement.SupplementService) *LambdaHandler {
	return &LambdaHandler{service: service}
}

func createSupplementService() *supplement.SupplementService {
	repo := createPostgresSupplementRepository()
	return supplement.NewSupplementService(repo)
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
