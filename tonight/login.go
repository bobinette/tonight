package tonight

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func RegisterLoginHandler(e *echo.Echo, jwtKey []byte, userRepository UserRepository) {
	h := loginHandler{
		userService: userService{
			jwtKey: jwtKey,
			repo:   userRepository,
		},
	}

	e.POST("/api/login", h.login)
	e.POST("/api/logout", h.logout)
}

type loginHandler struct {
	userService userService
}

func (*loginHandler) loginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func (h *loginHandler) login(c echo.Context) error {
	defer c.Request().Body.Close()

	var body struct {
		UserName string `json:"username"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	if body.UserName == "" {
		return errors.New("invalid empty username")
	}

	ctx := c.Request().Context()

	user, err := h.userService.getOrCreate(ctx, body.UserName)
	if err != nil {
		return err
	}

	accessToken, err := h.userService.token(ctx, user)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   accessToken,
		Expires: time.Now().Add(60 * 24 * time.Hour),
	})

	return c.NoContent(http.StatusNoContent)
}

func (*loginHandler) logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return c.NoContent(http.StatusNoContent)
}
