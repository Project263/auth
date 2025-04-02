package handlers

import (
	"auth/internal/config"
	"auth/internal/jwt"
	"auth/internal/services"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service *services.UserService
	cfg     *config.Config
}

func NewUserHandler(service *services.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{service: service, cfg: cfg}
}

func (h *UserHandler) Login(e echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := new(LoginRequest)

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Плохой запрос"})
	}
	ctx := e.Request().Context()

	user, err := h.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверные данные"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверные данные"})
	}

	token, err := jwt.GenerateJWT(h.cfg, user.Id, user.Role)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Ошибка создания токена"})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

func (h *UserHandler) Me(e echo.Context) error {
	tokenString := jwt.GetToken(e)
	logrus.Info(tokenString)
	if tokenString == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Нет токена"})
	}

	token, err := jwt.ParseJWT(h.cfg, tokenString)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, echo.Map{"error": "токен инвалидов"})
	}

	userId, ok := token["user_id"].(string)
	logrus.Info(userId, ok)
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

	return e.JSON(http.StatusOK, echo.Map{"message": "Выходим выходим"})
}
