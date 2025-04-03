package handlers

import (
	"auth/internal/config"
	"auth/internal/jwt"
	"auth/internal/services"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *services.UserService
	cfg     *config.Config
}

func NewUserHandler(service *services.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{service: service, cfg: cfg}
}

func (h *UserHandler) Me(e echo.Context) error {
	tokenString := jwt.GetToken(e)
	if tokenString == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Нет токена"})
	}

	claims, err := jwt.ParseJWT(h.cfg, tokenString)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, echo.Map{"error": "токен инвалидов"})
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return e.JSON(http.StatusUnauthorized, echo.Map{"error": "токен инвалидов"})
	}

	user, err := h.service.GetUserById(context.Background(), userId)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, echo.Map{"error": "токен инвалидов"})
	}

	return e.JSON(http.StatusOK, user)
}

func (h *UserHandler) Logout(e echo.Context) error {
	e.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "1",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   h.cfg.DOMAIN,
		MaxAge:   -1,
	})

	return e.Redirect(http.StatusFound, h.cfg.FRONT_URL)
}
