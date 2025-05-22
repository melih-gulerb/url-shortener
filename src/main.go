package main

import (
	"url-shortener/src/configs"
	"url-shortener/src/handlers"
	"url-shortener/src/helpers"
	"url-shortener/src/repositories"
)

func main() {
	helpers.LoadEnvFile(".env")
	cfg := helpers.LoadConfigFromEnv()

	mongo := configs.InitializeMongo(cfg.MongoURI)
	echo := configs.InitializeEcho()

	urlRepository := repositories.NewURLRepository(mongo.Database("urls"), "urls")
	urlHandler := handlers.NewURLHandler(urlRepository)

	echo.POST("/shorten", urlHandler.CreateShortURL)
}
