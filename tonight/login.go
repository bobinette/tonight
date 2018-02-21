package tonight

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

const userInfoScope = "https://www.googleapis.com/auth/userinfo.email"

func RegisterLoginHandler(e *echo.Echo, jwtKey []byte, cookieSecret []byte, clientID, clientSecret, redirectURL string, userRepository UserRepository) {
	h := loginHandler{
		userService: userService{
			jwtKey: jwtKey,
			repo:   userRepository,
		},

		sessionStore: sessions.NewCookieStore(cookieSecret),
		oauth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{userInfoScope},
			Endpoint:     google.Endpoint,
		},
	}

	e.GET("/api/oauth2/login", h.loginURL)
	e.GET("/api/oauth2/callback", h.oauth2Callback)
	e.POST("/api/logout", h.logout)
}

type loginHandler struct {
	userService userService

	sessionStore sessions.Store
	oauth2Config *oauth2.Config
}

// From:
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/getting-started/bookshelf/app/auth.go
func (h *loginHandler) loginURL(c echo.Context) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	sessionID := id.String()

	r := c.Request()
	oauthFlowSession, err := h.sessionStore.New(r, sessionID)
	if err != nil {
		return err
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	if err := oauthFlowSession.Save(r, c.Response()); err != nil {
		return err
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	url := h.oauth2Config.AuthCodeURL(sessionID, oauth2.AccessTypeOffline)
	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

// From:
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/getting-started/bookshelf/app/auth.go
func (h *loginHandler) oauth2Callback(c echo.Context) error {
	r := c.Request()

	state := r.FormValue("state")
	oauthFlowSession, err := h.sessionStore.Get(r, state)
	if err != nil {
		return err
	}
	if oauthFlowSession.Name() != state || state == "" {
		return errors.New("invalid state")
	}

	// Remove the session cookie
	c.SetCookie(&http.Cookie{
		Name:    oauthFlowSession.Name(),
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	code := r.FormValue("code")
	tok, err := h.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	client := h.oauth2Config.Client(c.Request().Context(), tok)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}

	ctx := c.Request().Context()
	user, err := h.userService.getOrCreate(ctx, body.Email)
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
		Path:    "/",
	})

	return c.Redirect(http.StatusFound, "http://127.0.0.1:9090/")
}

func (*loginHandler) logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return c.NoContent(http.StatusNoContent)
}
