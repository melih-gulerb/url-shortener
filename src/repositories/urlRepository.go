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

func (r *URLRepository) GetShortCodeByOriginalURL(ctx context.Context, originalURLValue string) (string, error) {
	filter := bson.M{"originalUrl": originalURLValue}
	opts := options.FindOne().SetProjection(bson.M{"shortCode": 1})
	var result bson.M
	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if url, ok := result["shortCode"].(string); ok {
		return url, nil
	}

	return "", errors.New("originalUrl not found in document or not a string")
}

func (r *URLRepository) GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	filter := bson.M{"shortCode": shortCode}
	opts := options.FindOne().SetProjection(bson.M{"originalUrl": 1})
	var result bson.M
	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if url, ok := result["originalUrl"].(string); ok {
		return url, nil
	}

	return "", errors.New("originalUrl not found in document or not a string")
}

func (r *URLRepository) Insert(ctx context.Context, urlData *models.UrlCollection) (*mongo.InsertOneResult, error) {
	now := time.Now()
	urlData.CreatedAt = now
	urlData.UpdatedAt = now

	return r.collection.InsertOne(ctx, urlData)
}
