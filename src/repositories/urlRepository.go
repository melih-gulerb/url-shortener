package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"url-shortener/src/models"
)

type URLRepository struct {
	collection *mongo.Collection
}

func NewURLRepository(db *mongo.Database, collectionName string) *URLRepository {
	return &URLRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *URLRepository) IsExisting(ctx context.Context, originalURL, shortCode string) (bool, error) {
	var orConditions []bson.M

	if originalURL != "" {
		orConditions = append(orConditions, bson.M{"originalUrl": originalURL})
	}
	if shortCode != "" {
		orConditions = append(orConditions, bson.M{"shortCode": shortCode})
	}

	filter := bson.M{"$or": orConditions}
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})
	var result bson.M
	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *URLRepository) Insert(ctx context.Context, urlData *models.UrlCollection) (*mongo.InsertOneResult, error) {
	now := time.Now()
	urlData.CreatedAt = now
	urlData.UpdatedAt = now

	return r.collection.InsertOne(ctx, urlData)
}
