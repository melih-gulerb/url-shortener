package models

import "time"

type CreateShortURLRequest struct {
	OriginalUrl string     `json:"originalUrl" validate:"required,url"`
	ShortCode   string     `json:"shortCode,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type GetOriginalURLRequest struct {
	ShortCode string `json:"shortCode" validate:"required"`
}
