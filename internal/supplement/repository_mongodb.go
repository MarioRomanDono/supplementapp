package supplement

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDBSupplementRepository struct {
	collection *mongo.Collection
}

func NewMongoDBSupplementRepository(collection *mongo.Collection) *MongoDBSupplementRepository {
	return &MongoDBSupplementRepository{collection: collection}
}

func (r *MongoDBSupplementRepository) FindByGtin(ctx context.Context, gtin string) (*Supplement, error) {
	var supplement Supplement
	err := r.collection.FindOne(ctx, bson.D{{Key: "gtin", Value: gtin}}).Decode(&supplement)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &supplement, err

}

func (r *MongoDBSupplementRepository) Create(ctx context.Context, supplement Supplement) error {
	_, err := r.collection.InsertOne(ctx, supplement)
	return err
}

func (r *MongoDBSupplementRepository) Update(ctx context.Context, supplement Supplement) error {
	_, err := r.collection.ReplaceOne(ctx, bson.D{{Key: "gtin", Value: supplement.Gtin}}, supplement)
	return err
}

func (r *MongoDBSupplementRepository) Delete(ctx context.Context, supplement Supplement) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{Key: "gtin", Value: supplement.Gtin}})
	return err
}

func (r *MongoDBSupplementRepository) ListAll(ctx context.Context) ([]Supplement, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	supplements := []Supplement{}
	if err := cursor.All(ctx, &supplements); err != nil {
		return nil, err
	}

	return supplements, nil
}