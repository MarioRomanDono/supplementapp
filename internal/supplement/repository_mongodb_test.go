package supplement_test

import (
	"context"
	"os"
	"testing"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setup(t *testing.T, ctx context.Context) (*mongo.Collection, *mongo.Client) {
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

func teardown(t *testing.T, ctx context.Context, client *mongo.Client) {
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

func TestMongoDBSupplementRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("success", func(t *testing.T) {
		context := context.Background()
		collection, client := setup(t, context)
		teardown(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		want := supplement.Supplement{
			Gtin: "1234567890123",
			Name: "name",
			Brand: "brand",
			Flavor: "flavor",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		err := repo.Create(context, want)

		if err != nil {
			t.Errorf("MongoDBSupplementRepository.Create() error = %v, want nil", err)
		}

		var got supplement.Supplement
		collection.FindOne(context, bson.D{{Key: "gtin", Value: want.Gtin}}).Decode(&got)
		
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("MongoDBSupplementRepository.Create() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestMongoDBSupplementRepository_FindByGtin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("not found", func(t *testing.T) {
		context := context.Background()
		collection, client := setup(t, context)
		teardown(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		got, err := repo.FindByGtin(context, "1234567890123")

		if err != nil {
			t.Errorf("MongoDBSupplementRepository.FindByGtin() error = %v, want nil", err)
		}

		if got != nil {
			t.Errorf("MongoDBSupplementRepository.FindByGtin() got = %v, want nil", got)
		}
	})

	t.Run("success", func(t *testing.T) {
		context := context.Background()
		collection, client := setup(t, context)
		teardown(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		want := &supplement.Supplement{
			Gtin: "1234567890123",
			Name: "name",
			Brand: "brand",
			Flavor: "flavor",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		collection.InsertOne(context, want)

		got, err := repo.FindByGtin(context, want.Gtin)

		if err != nil {
			t.Errorf("MongoDBSupplementRepository.FindByGtin() error = %v, want nil", err)
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("MongoDBSupplementRepository.FindByGtin() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestMongoDBSupplementRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("success", func(t *testing.T) {
		context := context.Background()
		collection, client := setup(t, context)
		teardown(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		want := supplement.Supplement{
			Gtin: "1234567890123",
			Name: "name",
			Brand: "brand",
			Flavor: "flavor",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		collection.InsertOne(context, want)

		want.Name = "new name"
		err := repo.Update(context, want)

		if err != nil {
			t.Errorf("MongoDBSupplementRepository.Update() error = %v, want nil", err)
		}

		var got supplement.Supplement
		collection.FindOne(context, bson.D{{Key: "gtin", Value: want.Gtin}}).Decode(&got)
		
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("MongoDBSupplementRepository.Update() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestMongoDBSupplementRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("success", func(t *testing.T) {
		context := context.Background()
		collection, client := setup(t, context)
		teardown(t, context, client)

		repo := supplement.NewMongoDBSupplementRepository(collection)
		s := supplement.Supplement{
			Gtin: "1234567890123",
			Name: "name",
			Brand: "brand",
			Flavor: "flavor",
			Carbohydrates: 1.0,
			Electrolytes: 1.0,
			Maltodextrose: 1.0,
			Fructose: 1.0,
			Caffeine: 1.0,
			Sodium: 1.0,
			Protein: 1.0,
		}
		collection.InsertOne(context, s)

		err := repo.Delete(context, s)

		if err != nil {
			t.Errorf("MongoDBSupplementRepository.Delete() error = %v, want nil", err)
		}

		var got *supplement.Supplement
		collection.FindOne(context, bson.D{{Key: "gtin", Value: s.Gtin}}).Decode(&got)
		
		if got != nil {
			t.Errorf("MongoDBSupplementRepository.Delete() got = %v, want nil", got)
		}
	})
}