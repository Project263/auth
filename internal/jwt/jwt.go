package jwt

import (
	"auth/internal/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(cfg *config.Config, userID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(cfg.SECRET))
}

func ParseJWT(cfg *config.Config, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.SECRET), nil
	})
}
