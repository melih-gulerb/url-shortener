package main

import (
	"github.com/labstack/echo/v4/middleware"
	"url-shortener/src/configs"
	"url-shortener/src/handlers"
	"url-shortener/src/helpers"
	"url-shortener/src/middlewares"
	"url-shortener/src/repositories"
)

func main() {
	helpers.LoadEnvFile(".env")
	cfg := helpers.LoadConfigFromEnv()

	mongo := configs.InitializeMongo(cfg.MongoURI)
	echo := configs.InitializeEcho()

	urlRepository := repositories.NewURLRepository(mongo.Database(cfg.MongoURLShortenerDatabase), cfg.MongoURLShortenerCollection)
	urlHandler := handlers.NewURLHandler(urlRepository, cfg.BaseURL)

	echo.Use(middleware.Recover())
	echo.Use(middlewares.ResponseBodyLogger)
	echo.POST("/shorten", urlHandler.CreateShortURL)
	echo.POST("/original", urlHandler.GetOriginalURL)

	echo.Logger.Fatal(echo.Start(":1323"))
}
