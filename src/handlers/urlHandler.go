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
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Data:    nil,
			Message: models.Messages.FailedToCreateShortCode,
		})
	}

	if req.OriginalUrl == "" {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Data:    nil,
			Message: models.Messages.OriginalURLCannotBeEmpty,
		})
	}

	var shortCodeToUse string
	if req.ShortCode != "" {
		shortCodeToUse = req.ShortCode
	} else {
		shortCodeToUse = helpers.GenerateShortCode(req.OriginalUrl)
		if shortCodeToUse == "" {
			return c.JSON(http.StatusInternalServerError, "Failed to generate short code")
		}
	}

	urlDoc := &models.UrlCollection{
		OriginalUrl: req.OriginalUrl,
		ShortCode:   shortCodeToUse,
		AccessCount: 0,
	}

	shortCode, err := h.repo.GetOriginalURLByShortCode(c.Request().Context(), shortCodeToUse)
	if shortCode != "" && req.ShortCode != "" {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Data:    nil,
			Message: models.Messages.ShortCodeAlreadyExists,
		})
	}

	shortCode, err = h.repo.GetShortCodeByOriginalURL(c.Request().Context(), req.OriginalUrl)
	if shortCode != "" {
		return c.JSON(http.StatusOK, models.BaseResponse{
			Data:    shortCode,
			Message: models.Messages.Success,
		})
	}

	_, err = h.repo.Insert(c.Request().Context(), urlDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Data:    nil,
			Message: models.Messages.FailedToCreateShortCode,
		})
	}

	return c.JSON(http.StatusCreated, models.BaseResponse{
		Data:    shortCode,
		Message: models.Messages.Success,
	})
}

func (h *URLHandler) GetOriginalURL(c echo.Context) error {
	req := new(models.GetOriginalURLRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: models.Messages.FailedToGetOriginalURL,
			Data:    nil,
		})
	}

	if req.ShortCode == "" {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: models.Messages.ShortCodeCannotBeEmpty,
			Data:    nil,
		})
	}

	originalUrl, err := h.repo.GetOriginalURLByShortCode(c.Request().Context(), req.ShortCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: models.Messages.FailedToGetOriginalURL,
			Data:    nil,
		})
	}

	if originalUrl == "" {
		return c.JSON(http.StatusNotFound, models.BaseResponse{
			Message: models.Messages.OriginalURLNotFound,
		})
	}

	return c.JSON(http.StatusOK, models.BaseResponse{
		Data:    originalUrl,
		Message: models.Messages.Success,
	})
}
