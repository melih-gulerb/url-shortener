package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/melih-gulerb/go-logger/logging"
)

func ResponseBodyLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		logging.Default().Info(req.Method, req.URL.Path, res.Status)

		err := next(c)

		return err
	}
}
