package auth

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

func DevMiddleware(id, name string, s Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := User{
				ID:   id,
				Name: name,
			}
			ctx := c.Request().Context()
			if err := s.Ensure(ctx, &user); err != nil {
				return fmt.Errorf("ensuring user: %w", err)
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func CaddyMiddleware(s Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("Token-Claim-Sub")
			if id == "" {
				return errors.New("missing user") // TODO: 401
			}

			user := User{
				ID:   id,
				Name: c.Request().Header.Get("Token-Claim-Name"),
			}

			ctx := c.Request().Context()
			if err := s.Ensure(ctx, &user); err != nil {
				return fmt.Errorf("error ensuring user: %w", err)
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
