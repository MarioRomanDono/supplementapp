package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marioromandono/supplementapp/cmd/http-server"
	"github.com/marioromandono/supplementapp/internal/supplement"
	"github.com/marioromandono/supplementapp/internal/supplement/persistence/postgres"

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
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}

	err = container.Snapshot(ctx, tcpostgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Panic(err)
	}

	dbUrl, err = container.ConnectionString(ctx)
	if err != nil {
		log.Panic(err)
	}

	m.Run()
}

func TestGetSupplement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		gtin := "123"
		request := httptest.NewRequest("GET", "/supplement/"+gtin, nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusNotFound
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: fmt.Sprintf("%s: %s", gtin, supplement.ErrNotFound),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("existing", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		want := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, want)

		request := httptest.NewRequest("GET", "/supplement/"+want.Gtin, nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusOK
		wantBodyJSON, _ := json.Marshal(want)
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})
}

func TestCreateSupplement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("without body", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		request := httptest.NewRequest("POST", "/supplement", nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusBadRequest
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: io.EOF.Error(),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("invalid json", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		body := []byte(`{"gtin": "1234567890123"]`)
		request := httptest.NewRequest("POST", "/supplement", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		err := json.Unmarshal(body, &supplement.Supplement{})
		wantCode := http.StatusBadRequest
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: err.Error(),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("invalid supplement", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		s := &supplement.Supplement{
			Gtin:          "1234567890123",
			Name:          "name",
			Brand:         "brand",
			Flavor:        "flavor",
			Carbohydrates: -1.0,
		}
		body, _ := json.Marshal(s)
		request := httptest.NewRequest("POST", "/supplement", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		wantCode := http.StatusBadRequest
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: fmt.Sprintf("%s: carbohydrates %f is invalid, it must be greater or equal to zero", supplement.ErrInvalidSupplement, s.Carbohydrates),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("already exists", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		s := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, s)

		body, _ := json.Marshal(s)
		request := httptest.NewRequest("POST", "/supplement", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		wantCode := http.StatusConflict
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: fmt.Sprintf("%s: %s", s.Gtin, supplement.ErrAlreadyExists),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("created", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		s := &supplement.Supplement{
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
		}
		body, _ := json.Marshal(s)
		request := httptest.NewRequest("POST", "/supplement", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		wantCode := http.StatusCreated

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), "")
		assertHeader(t, response.Header(), "Location", "/supplement/"+s.Gtin)
	})
}

func TestUpdateSupplement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("without body", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		gtin := "123"
		request := httptest.NewRequest("PUT", "/supplement/"+gtin, nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusBadRequest
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: io.EOF.Error(),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("invalid json", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		body := []byte(`{"gtin": "1234567890123"]`)
		request := httptest.NewRequest("PUT", "/supplement/1234567890123", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		err := json.Unmarshal(body, &supplement.Supplement{})
		wantCode := http.StatusBadRequest
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: err.Error(),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		gtin := "123"
		body, _ := json.Marshal(&supplement.Supplement{
			Gtin: gtin,
			Name: "Test",
		})
		request := httptest.NewRequest("PUT", "/supplement/"+gtin, bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		wantCode := http.StatusNotFound
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: fmt.Sprintf("%s: %s", gtin, supplement.ErrNotFound),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("invalid updatable supplement", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		s := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, s)

		s.Carbohydrates = -1.0
		body, _ := json.Marshal(s)
		request := httptest.NewRequest("PUT", "/supplement/"+s.Gtin, bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		wantCode := http.StatusBadRequest
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: fmt.Sprintf("%s: carbohydrates %f is invalid, it must be greater or equal to zero", supplement.ErrInvalidSupplement, s.Carbohydrates),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("updated", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		s := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, s)

		s.Name = "Updated"
		body, _ := json.Marshal(s)
		request := httptest.NewRequest("PUT", "/supplement/"+s.Gtin, bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		wantCode := http.StatusOK

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), "")
	})
}

func TestDeleteSupplement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		gtin := "123"
		request := httptest.NewRequest("DELETE", "/supplement/"+gtin, nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusNotFound
		wantBodyJSON, _ := json.Marshal(&main.ErrorResponseBody{
			Code:    wantCode,
			Message: fmt.Sprintf("%s: %s", gtin, supplement.ErrNotFound),
		})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
	})

	t.Run("deleted", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		s := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, s)

		request := httptest.NewRequest("DELETE", "/supplement/"+s.Gtin, nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusNoContent

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), "")
	})
}

func TestListAllSupplements(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		dbPool := getPool(t, ctx)
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		request := httptest.NewRequest("GET", "/supplement", nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusOK
		wantBodyJSON, _ := json.Marshal([]supplement.Supplement{})
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
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
		server := main.NewServer(supplement.NewSupplementService(postgres.NewSupplementRepository(dbPool)))

		want := []supplement.Supplement{
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
		for _, s := range want {
			insertSupplement(t, ctx, dbPool, s)
		}

		request := httptest.NewRequest("GET", "/supplement", nil)
		response := httptest.NewRecorder()
		wantCode := http.StatusOK
		wantBodyJSON, _ := json.Marshal(want)
		wantBody := string(wantBodyJSON) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, wantCode)
		assertResponseBody(t, response.Body.String(), wantBody)
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

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect response body, got %s, want %s", got, want)
	}
}

func assertHeader(t testing.TB, got http.Header, key, want string) {
	t.Helper()
	if got.Get(key) != want {
		t.Errorf("incorrect header, got %s, want %s", got.Get(key), want)
	}
}
