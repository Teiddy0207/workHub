package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"

	"workHub/constant"
	"workHub/logger"
	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(ctx context.Context, publicKey *rsa.PublicKey, tokenStr string) (JwtClaim, error) {
	var (
		claims JwtClaim
		err    error
	)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, keyFunc)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			logger.Error("jwt", "VerifyToken", fmt.Sprintf("JWT validation error: %v", ve))
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				logger.Warn("jwt", "VerifyToken", "Token expired")
				return claims, constant.ErrUnAuthentication
			}
		} else {
			logger.Error("jwt", "VerifyToken", fmt.Sprintf("JWT parse error: %v", err))
		}
		return claims, constant.ErrUnAuthentication
	}

	if !token.Valid {
		logger.Warn("jwt", "VerifyToken", "Token invalid")
		return claims, constant.ErrUnAuthentication
	}

	return claims, nil
}
