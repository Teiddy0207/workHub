package helper

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func GenerateResetToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"type":  "reset_password",
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
