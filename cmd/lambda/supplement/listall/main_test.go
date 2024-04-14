package main_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/marioromandono/supplementapp/cmd/lambda/supplement/listall"
	"github.com/marioromandono/supplementapp/internal/supplement"
	"github.com/marioromandono/supplementapp/internal/supplement/persistence/postgres"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var container *tcpostgres.PostgresContainer
var dbUrl string

const tableName string = "Supplements"

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "supplementapp"
	dbUser := "postgres"
	dbPassword := "password"

	var err error

	container, err = tcpostgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	_, _, err = container.Exec(ctx, []string{
		"psql", "-U", dbUser, "-d", dbName, "-c",
		"CREATE TABLE " + tableName + " ( " +
			"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, " +
			"gtin VARCHAR UNIQUE, " +
			"name VARCHAR, " +
			"brand VARCHAR, " +
			"flavor VARCHAR, " +
			"carbohydrates REAL, " +
			"electrolytes REAL, " +
			"maltodextrose REAL, " +
			"fructose REAL, " +
			"caffeine REAL, " +
			"sodium REAL, " +
			"protein REAL " +
			")",
	})
	if err != nil {
		panic(err)
	}

	err = container.Snapshot(ctx, tcpostgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		panic(err)
	}

	dbUrl, err = container.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestLambdaHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Setenv("POSTGRES_URL", dbUrl)

	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)

		want := events.APIGatewayV2HTTPResponse{
			Body:       "[]",
			StatusCode: 200,
		}

		handler := main.NewLambdaHandler(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))
		got, err := handler.Handle(ctx, events.APIGatewayV2HTTPRequest{})

		if err != nil {
			t.Errorf("LambdaHandler() error = %v, want nil", err)
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("LambdaHandler() error (-got +want):\n%s", diff)
		}
	})

	t.Run("not empty", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)

		ss := []supplement.Supplement{
			{
				Gtin:          "1234567890123",
				Name:          "Test",
				Brand:         "Test",
				Flavor:        "Test",
				Carbohydrates: 1.0,
				Electrolytes:  1.0,
				Maltodextrose: 1.0,
				Fructose:      1.0,
				Caffeine:      1.0,
				Sodium:        1.0,
				Protein:       1.0,
			},
			{
				Gtin:          "1234567890124",
				Name:          "Test",
				Brand:         "Test",
				Flavor:        "Test",
				Carbohydrates: 1.0,
				Electrolytes:  1.0,
				Maltodextrose: 1.0,
				Fructose:      1.0,
				Caffeine:      1.0,
				Sodium:        1.0,
				Protein:       1.0,
			},
		}
		for _, s := range ss {
			insertSupplement(t, ctx, dbPool, s)
		}

		ssJson, _ := json.Marshal(ss)
		want := events.APIGatewayV2HTTPResponse{
			Body:       string(ssJson),
			StatusCode: 200,
		}

		handler := main.NewLambdaHandler(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))
		got, err := handler.Handle(ctx, events.APIGatewayV2HTTPRequest{})

		if err != nil {
			t.Errorf("LambdaHandler() error = %v, want nil", err)
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("LambdaHandler() error (-got +want):\n%s", diff)
		}
	})
}

func getPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		dbPool.Close()
	})

	return dbPool
}

func insertSupplement(t *testing.T, ctx context.Context, dbPool *pgxpool.Pool, s supplement.Supplement) {
	t.Helper()
	_, err := dbPool.Exec(
		ctx,
		"INSERT INTO "+tableName+
			" (gtin, name, brand, flavor, carbohydrates, electrolytes, maltodextrose, fructose, caffeine, sodium, protein) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		s.Gtin, s.Name, s.Brand, s.Flavor, s.Carbohydrates, s.Electrolytes, s.Maltodextrose, s.Fructose, s.Caffeine, s.Sodium, s.Protein,
	)

	if err != nil {
		t.Fatal(err)
	}
}
