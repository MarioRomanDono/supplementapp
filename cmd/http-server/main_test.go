package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/marioromandono/supplementapp/cmd/http-server"
	"github.com/marioromandono/supplementapp/internal/supplement"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetSupplement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("not found", func(t *testing.T) {
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		want := &supplement.Supplement{
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
		_, err := collection.InsertOne(context, want)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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

	t.Run("already exists", func(t *testing.T) {
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		_, err := collection.InsertOne(context, s)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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

	t.Run("updated", func(t *testing.T) {
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		_, err := collection.InsertOne(context, s)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		_, err := collection.InsertOne(context, s)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
		context := context.Background()
		collection := setup(t, context)
		server := main.NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

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
			_, err := collection.InsertOne(context, s)
			if err != nil {
				t.Fatalf("could not insert document: %v", err)
			}
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

func setup(t *testing.T, ctx context.Context) *mongo.Collection {
	t.Helper()
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("could not load .env file: %v", err)
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		t.Fatalf("could not connect to mongodb: %v", err)
	}

	dbName := os.Getenv("MONGODB_DB")
	collectionName := "supplements"

	t.Cleanup(func() {
		err := client.Database(dbName).Collection(collectionName).Drop(ctx)
		if err != nil {
			t.Fatalf("could not drop collection: %v", err)
		}

		err = client.Disconnect(ctx)
		if err != nil {
			t.Fatalf("could not disconnect from mongodb: %v", err)
		}
	})

	return client.Database(dbName).Collection(collectionName)
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
