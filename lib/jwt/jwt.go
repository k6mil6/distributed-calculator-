package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"time"
)

func NewToken(user model.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	return token.SignedString([]byte(secret))
}
