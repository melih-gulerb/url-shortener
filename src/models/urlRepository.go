package models

import "time"

type UrlCollection struct {
	ID          string     `bson:"_id"`
	OriginalUrl string     `bson:"originalUrl"`
	ShortCode   string     `bson:"shortCode"`
	CreatedAt   time.Time  `bson:"createdAt"`
	UpdatedAt   time.Time  `bson:"updatedAt"`
	DeletedAt   *time.Time `bson:"deletedAt"`
	ExpiresAt   *time.Time `bson:"expiresAt"`
	AccessCount int        `bson:"accessCount"`
}
