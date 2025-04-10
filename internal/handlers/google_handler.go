package handlers

import (
	"auth/internal/config"
	"auth/internal/jwt"
	"auth/internal/models"
	"auth/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleHandler struct {
	cfg     *config.Config
	service *services.GoogleService
	gConfig *oauth2.Config
}

func NewGoogleHandler(service *services.GoogleService, cfg *config.Config) *GoogleHandler {
	var googleOauthConfig = &oauth2.Config{
		ClientID:     cfg.GOOGLE_CLIENT_ID,
		ClientSecret: cfg.GOOGLE_SECRET,
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", cfg.SSO_URL),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleHandler{cfg: cfg, service: service, gConfig: googleOauthConfig}
}

func (h *GoogleHandler) HandleGoogleLogin(c echo.Context) error {
	url := h.gConfig.AuthCodeURL("random-state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *GoogleHandler) HandleGoogleCallback(e echo.Context) error {
	code := e.QueryParam("code")
	ctx := e.Request().Context()

	token, err := h.gConfig.Exchange(context.Background(), code)
	if err != nil {
		logrus.Error("Ошибка при обмене кода:", err)
		return e.String(http.StatusInternalServerError, "Ошибка при авторизации")
	}

	client := h.gConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logrus.Error("Ошибка при получении данных пользователя:", err)
		return e.String(http.StatusInternalServerError, "Ошибка при авторизации")
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	email := userInfo["email"].(string)
	username := userInfo["name"].(string)
	avatar := userInfo["picture"].(string)

	userData := models.User{
		Email:    email,
		Username: username,
		Avatar:   avatar,
	}

	userId, err := h.service.CreateUserByGoogle(ctx, userData)
	logrus.Info("handler", userId, err)
	if err != nil {
		logrus.Error(err)
		return e.String(http.StatusInternalServerError, "Ошибка при авторизации")
	}

	jwtToken, err := jwt.GenerateJWT(h.cfg, userId, "user")
	if err != nil {
		logrus.Error(err)
		return e.String(http.StatusInternalServerError, "Ошибка создания токена")
	}

	e.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   h.cfg.DOMAIN,
	})

	return e.Redirect(http.StatusFound, h.cfg.FRONT_URL)
}
