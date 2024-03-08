package supplement_test

import (
	"context"
	"os"
	"testing"
	"reflect"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
)

func setupTest(t *testing.T, ctx context.Context) (*mongo.Collection, *mongo.Client) {
	t.Helper()
	err := godotenv.Load("../../.env.offline")
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

func TestCreate(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		expected := newRandomSupplement()
		err := repo.Create(context, *expected)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		var actual *supplement.Supplement
		collection.FindOne(context, bson.D{{Key: "gtin", Value: expected.Gtin}}).Decode(&actual)
		
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

func TestFindByGtin(t *testing.T) {
	t.Run("successful find", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		expected := newRandomSupplement()
		collection.InsertOne(context, expected)

		actual, err := repo.FindByGtin(context, expected.Gtin)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("not found", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		actual, err := repo.FindByGtin(context, "1234567890123")

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if actual != nil {
			t.Errorf("expected nil, got %v", actual)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		expected := newRandomSupplement()
		collection.InsertOne(context, expected)

		expected.Name = "new name"
		err := repo.Update(context, *expected)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		var actual *supplement.Supplement
		collection.FindOne(context, bson.D{{Key: "gtin", Value: expected.Gtin}}).Decode(&actual)
		
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		context := context.Background()
		collection, client := setupTest(t, context)
		teardownTest(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		expected := newRandomSupplement()
		collection.InsertOne(context, expected)

		err := repo.Delete(context, *expected)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		var actual *supplement.Supplement
		collection.FindOne(context, bson.D{{Key: "gtin", Value: expected.Gtin}}).Decode(&actual)
		
		if actual != nil {
			t.Errorf("expected nil, got %v", actual)
		}
	})
}