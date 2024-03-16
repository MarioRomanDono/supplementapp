package main

import (
	"bytes"
	"context"
	"fmt"
	"encoding/json"
	"os"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetSupplement(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		gtin := "123"
		request := httptest.NewRequest("GET", "/supplement/"+gtin, nil)
		response := httptest.NewRecorder()
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(`{"code":%d,"message":"%s"}`+"\n", expectedCode, fmt.Sprintf("%s: %s", gtin, supplement.ErrNotFound))	

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), expectedBody)
	})

	t.Run("existing", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		expected := &supplement.Supplement{
			Gtin: "1234567890123",
			Name: "Test",
			Brand: "Test",
			Flavor: "Test",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		_, err := collection.InsertOne(context, expected)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

		request := httptest.NewRequest("GET", "/supplement/"+expected.Gtin, nil)
		response := httptest.NewRecorder()
		expectedCode := http.StatusOK
		json, _ := json.Marshal(expected)
		expectedBody := string(json) + "\n"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), expectedBody)
	})
}

func TestCreateSupplement(t *testing.T) {
	t.Run("already exists", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		expected := &supplement.Supplement{
			Gtin: "1234567890123",
			Name: "Test",
			Brand: "Test",
			Flavor: "Test",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		_, err := collection.InsertOne(context, expected)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

		body, _ := json.Marshal(expected)
		request := httptest.NewRequest("POST", "/supplement", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		expectedCode := http.StatusConflict
		expectedBody := fmt.Sprintf(`{"code":%d,"message":"%s"}`+"\n", expectedCode, fmt.Sprintf("%v: %s", *expected, supplement.ErrAlreadyExists))

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), expectedBody)
	})

	t.Run("created", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		expected := &supplement.Supplement{
			Gtin: "1234567890123",
			Name: "Test",
			Brand: "Test",
			Flavor: "Test",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		body, _ := json.Marshal(expected)
		request := httptest.NewRequest("POST", "/supplement", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		expectedCode := http.StatusCreated

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), "")
		assertHeader(t, response.Header(), "Location", "/supplement/"+expected.Gtin)
	})
}

func TestUpdateSupplement(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		gtin := "123"
		body, _ := json.Marshal(&supplement.Supplement{
			Gtin: gtin,
			Name: "Test",
		})
		request := httptest.NewRequest("PUT", "/supplement/"+gtin, bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(`{"code":%d,"message":"%s"}`+"\n", expectedCode, fmt.Sprintf("%s: %s", gtin, supplement.ErrNotFound))

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), expectedBody)
	})

	t.Run("updated", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		expected := &supplement.Supplement{
			Gtin: "1234567890123",
			Name: "Test",
			Brand: "Test",
			Flavor: "Test",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		_, err := collection.InsertOne(context, expected)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

		expected.Name = "Updated"
		body, _ := json.Marshal(expected)
		request := httptest.NewRequest("PUT", "/supplement/"+expected.Gtin, bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		expectedCode := http.StatusOK

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), "")
	})
}

func TestDeleteSupplement(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		gtin := "123"
		request := httptest.NewRequest("DELETE", "/supplement/"+gtin, nil)
		response := httptest.NewRecorder()
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(`{"code":%d,"message":"%s"}`+"\n", expectedCode, fmt.Sprintf("%s: %s", gtin, supplement.ErrNotFound))

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), expectedBody)
	})

	t.Run("deleted", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)
		server := NewServer(supplement.NewSupplementService(supplement.NewMongoDBSupplementRepository(collection)))

		expected := &supplement.Supplement{
			Gtin: "1234567890123",
			Name: "Test",
			Brand: "Test",
			Flavor: "Test",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		_, err := collection.InsertOne(context, expected)
		if err != nil {
			t.Fatalf("could not insert document: %v", err)
		}

		request := httptest.NewRequest("DELETE", "/supplement/"+expected.Gtin, nil)
		response := httptest.NewRecorder()
		expectedCode := http.StatusNoContent

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, expectedCode)
		assertResponseBody(t, response.Body.String(), "")
	})
}

func setupTest(t *testing.T, ctx context.Context) (*mongo.Collection, *mongo.Client) {
	t.Helper()
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("could not load .env file: %v", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err != nil {
		t.Fatalf("could not connect to mongodb: %v", err)
	}

	return client.Database(os.Getenv("MONGODB_DB")).Collection("supplements"), client
}

func teardownTest(t *testing.T, ctx context.Context, client *mongo.Client) {
	t.Cleanup(func() {
		err := client.Database(os.Getenv("MONGODB_DB")).Collection("supplements").Drop(ctx)
		if err != nil {
			t.Fatalf("could not drop collection: %v", err)
		}

		err = client.Disconnect(ctx)
		if err != nil {
			t.Fatalf("could not disconnect from mongodb: %v", err)
		}
	})
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