package tonight

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if err == middleware.ErrJWTMissing {
		code = http.StatusUnauthorized
	}

	c.JSON(code, map[string]interface{}{"error": err.Error()})
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"data": "ok"})
}
