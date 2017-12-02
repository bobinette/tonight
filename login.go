package tonight

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type User struct {
	ID   uint
	Name string

	TaskIDs []uint
}

type UserRepository interface {
	Get(ctx context.Context, id uint) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Insert(ctx context.Context, user *User) error

	AddTaskToUser(ctx context.Context, userID uint, taskID uint) error
}

type LoginService struct {
	Key        []byte
	Repository UserRepository
}

func (s *LoginService) LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func (s *LoginService) Login(c echo.Context) error {
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

	user, err := s.Repository.GetByName(ctx, body.UserName)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		user = User{Name: body.UserName}
		if err := s.Repository.Insert(ctx, &user); err != nil {
			return err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tonightClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 2, 0).Unix(),
			Issuer:    "tonight",
		},
	})

	tokenStr, err := token.SignedString(s.Key)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   tokenStr,
		Expires: time.Now().Add(24 * time.Hour),
	})

	return c.NoContent(http.StatusNoContent)
}

func (s *LoginService) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return c.NoContent(http.StatusNoContent)
}
