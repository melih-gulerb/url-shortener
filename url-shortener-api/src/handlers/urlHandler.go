package handlers

import (
	"github.com/melih-gulerb/go-logger/logging"
	"net/http"
	"url-shortener/src/helpers"
	models2 "url-shortener/src/models"
	"url-shortener/src/repositories"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	repo    *repositories.URLRepository
	baseUrl string
}

func NewURLHandler(repo *repositories.URLRepository, baseUrl string) *URLHandler {
	return &URLHandler{
		repo:    repo,
		baseUrl: baseUrl,
	}
}

func (h *URLHandler) CreateShortURL(c echo.Context) error {
	req := new(models2.CreateShortURLRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models2.BaseResponse{
			Data:    nil,
			Message: models2.Messages.FailedToCreateShortCode,
		})
	}

	if req.OriginalUrl == "" {
		return c.JSON(http.StatusBadRequest, models2.BaseResponse{
			Data:    nil,
			Message: models2.Messages.OriginalURLCannotBeEmpty,
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

	urlDoc := &models2.UrlCollection{
		OriginalUrl: req.OriginalUrl,
		ShortCode:   shortCodeToUse,
		AccessCount: 0,
	}

	shortCode, err := h.repo.GetOriginalURLByShortCode(c.Request().Context(), shortCodeToUse)
	if shortCode != "" && req.ShortCode != "" {
		return c.JSON(http.StatusBadRequest, models2.BaseResponse{
			Data:    nil,
			Message: models2.Messages.ShortCodeAlreadyExists,
		})
	}

	shortCode, err = h.repo.GetShortCodeByOriginalURL(c.Request().Context(), req.OriginalUrl)
	if shortCode != "" {
		return c.JSON(http.StatusOK, models2.BaseResponse{
			Data:    h.baseUrl + "/" + shortCodeToUse,
			Message: models2.Messages.Success,
		})
	}

	_, err = h.repo.Insert(c.Request().Context(), urlDoc)
	if err != nil {
		logging.Default().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models2.BaseResponse{
			Data:    nil,
			Message: models2.Messages.FailedToCreateShortCode,
		})
	}

	return c.JSON(http.StatusCreated, models2.BaseResponse{
		Data:    h.baseUrl + "/" + shortCodeToUse,
		Message: models2.Messages.Success,
	})
}

func (h *URLHandler) GetOriginalURL(c echo.Context) error {
	req := new(models2.GetOriginalURLRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models2.BaseResponse{
			Message: models2.Messages.FailedToGetOriginalURL,
			Data:    nil,
		})
	}

	if req.ShortCode == "" {
		return c.JSON(http.StatusBadRequest, models2.BaseResponse{
			Message: models2.Messages.ShortCodeCannotBeEmpty,
			Data:    nil,
		})
	}

	originalUrl, err := h.repo.GetOriginalURLByShortCode(c.Request().Context(), req.ShortCode)
	if err != nil {
		logging.Default().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models2.BaseResponse{
			Message: models2.Messages.FailedToGetOriginalURL,
			Data:    nil,
		})
	}

	if originalUrl == "" {
		return c.JSON(http.StatusNotFound, models2.BaseResponse{
			Message: models2.Messages.OriginalURLNotFound,
		})
	}

	return c.JSON(http.StatusOK, models2.BaseResponse{
		Data:    originalUrl,
		Message: models2.Messages.Success,
	})
}
