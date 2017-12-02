package tonight

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type tonightClaims struct {
	ID uint `json:"user_id"`

	jwt.StandardClaims
}

func JWTMiddleware(key []byte) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &tonightClaims{},
		SigningKey:  key,
		ContextKey:  "access_token",
		TokenLookup: "cookie:access_token",
	})
}

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if err == middleware.ErrJWTMissing {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if code == http.StatusUnauthorized {
		c.Render(code, "login", nil)
		return
	}

	c.JSON(code, map[string]interface{}{"error": err.Error()})
}
