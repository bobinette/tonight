package tonight

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	ErrMissingToken            = errors.New("missing token")
	ErrMissingUser             = errors.New("missing user")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInvalidClaims           = errors.New("invalid claims")
)

type tonightClaims struct {
	UserID uint `json:"user_id"`

	jwt.StandardClaims
}

func UserMiddleware(userRepository UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			email := c.Request().Header.Get("X-Forwarded-Email")

			user, err := userRepository.GetByName(c.Request().Context(), email)
			if err != nil {
				return err
			}

			if user.Name == "" {
				// Not found, we insert the user
				user = User{
					Name:    email,
					IsAdmin: false,
				}
				if err := userRepository.Insert(c.Request().Context(), &user); err != nil {
					return err
				}
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := loadUser(c)
			if err != nil {
				return err
			}

			if !user.IsAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
			}

			return next(c)
		}
	}
}

func loadUser(c echo.Context) (User, error) {
	user, ok := c.Get("user").(User)
	if !ok {
		return User{}, ErrMissingUser
	}
	return user, nil
}

func reloadUser(c echo.Context, userRepository UserRepository) (User, error) {
	user, ok := c.Get("user").(User)
	if !ok {
		return User{}, ErrMissingUser
	}

	user, err := userRepository.Get(c.Request().Context(), user.ID)
	if err != nil {
		return User{}, err
	}

	c.Set("user", user)
	return user, nil
}

func checkPermission(c echo.Context, taskID uint) error {
	user, err := loadUser(c)
	if err != nil {
		return err
	}

	for _, tID := range user.TaskIDs {
		if tID == taskID {
			return nil
		}
	}

	return ErrInsufficientPermissions
}
