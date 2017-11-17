package tonight

import (
	"net/http"

	"github.com/labstack/echo"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"data": "ok"})
}
