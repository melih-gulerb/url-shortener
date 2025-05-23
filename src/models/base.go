package models

type BaseResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var Messages = struct {
	Success                  string
	ShortCodeAlreadyExists   string
	FailedToCreateShortCode  string
	OriginalURLCannotBeEmpty string
	ShortCodeCannotBeEmpty   string
	FailedToGetOriginalURL   string
	OriginalURLNotFound      string
}{
	Success:                  "success",
	ShortCodeAlreadyExists:   "short code already exists",
	FailedToCreateShortCode:  "failed to create short code",
	OriginalURLCannotBeEmpty: "original url cannot be empty",
	ShortCodeCannotBeEmpty:   "short code cannot be empty",
	FailedToGetOriginalURL:   "failed to get original url",
	OriginalURLNotFound:      "original url not found",
}
