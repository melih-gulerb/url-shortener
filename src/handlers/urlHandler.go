package handlers

import (
	"net/http"
	"url-shortener/src/helpers"
	"url-shortener/src/models"
	"url-shortener/src/repositories"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	repo *repositories.URLRepository
}

func NewURLHandler(repo *repositories.URLRepository) *URLHandler {
	return &URLHandler{
		repo: repo,
	}
}

func (h *URLHandler) CreateShortURL(c echo.Context) error {
	req := new(models.CreateShortURLRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.OriginalUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "OriginalUrl is required")
	}

	var shortCodeToUse string
	if req.ShortCode != "" {
		shortCodeToUse = req.ShortCode
	} else {
		shortCodeToUse = helpers.GenerateShortCode(req.OriginalUrl)
		if shortCodeToUse == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate short code")
		}
	}

	urlDoc := &models.UrlCollection{
		OriginalUrl: req.OriginalUrl,
		ShortCode:   shortCodeToUse,
		AccessCount: 0,
	}

	existing, err := h.repo.IsExisting(c.Request().Context(), req.OriginalUrl, shortCodeToUse)
	if existing {
		return c.JSON(http.StatusBadRequest, "Short URL already exists")
	}

	_, err = h.repo.Insert(c.Request().Context(), urlDoc)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create short URL: "+err.Error())
	}

	return c.JSON(http.StatusCreated, urlDoc)
}
