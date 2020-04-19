package auth

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func userFromHeader(c echo.Context) (User, error) {
	id := c.Request().Header.Get("Token-Claim-Sub")
	if id == "" {
		// TODO: use the 401 error code
		return User{}, errors.New("no user")
	}

	name := c.Request().Header.Get("Token-Claim-Name")

	return User{
		ID:   id,
		Name: name,
	}, nil
}

func ExtractUser(c echo.Context) (User, error) {
	v := c.Get("user")
	if v == nil {
		return User{}, errors.New("missing user") // TODO 401
	}

	user, ok := v.(User)
	if !ok {
		return User{}, errors.New("missing user") // TODO 401
	}

	return user, nil
}
