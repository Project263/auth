package handlers

import (
	"auth/internal/config"
	"auth/internal/jwt"
	"auth/internal/services"
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

func (s *UserHandler) Login(e echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := new(LoginRequest)

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Плохой запрос"})
	}
	ctx := e.Request().Context()

	user, err := s.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверные данные"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверные данные"})
	}

	token, err := jwt.GenerateJWT(s.cfg, user.Id, user.Role)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Ошибка создания токена"})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
}
