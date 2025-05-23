package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/melih-gulerb/go-logger/logging"
)

func ResponseBodyLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		logging.Default().Info(req.Method, req.URL)

		err := next(c)

		return err
	}
}
