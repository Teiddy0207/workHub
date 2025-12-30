package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"workHub/constant"
	"workHub/logger"

	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(ctx context.Context, publicKey *rsa.PublicKey, tokenStr string) (JwtClaim, error) {
	var (
		claims JwtClaim
		err    error
	)

	// Kiểm tra token không rỗng và có ít nhất một số ký tự
	if tokenStr == "" {
		logger.Error("jwt", "VerifyToken", "Empty token string")
		return claims, constant.ErrUnAuthentication
	}

	// Kiểm tra token có chứa ít nhất 2 dấu chấm (format JWT: header.payload.signature)
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		logger.Error("jwt", "VerifyToken", fmt.Sprintf("Invalid token format: expected 3 parts, got %d", len(parts)))
		return claims, constant.ErrUnAuthentication
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, keyFunc)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			logger.Error("jwt", "VerifyToken", fmt.Sprintf("JWT validation error: %v, errors: %d", ve, ve.Errors))
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				logger.Warn("jwt", "VerifyToken", "Token expired")
				return claims, constant.ErrUnAuthentication
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				logger.Error("jwt", "VerifyToken", "Token malformed - check token format")
				return claims, constant.ErrUnAuthentication
			}
		} else {
			logger.Error("jwt", "VerifyToken", fmt.Sprintf("JWT parse error: %v, token length: %d", err, len(tokenStr)))
		}
		return claims, constant.ErrUnAuthentication
	}

	if !token.Valid {
		logger.Warn("jwt", "VerifyToken", "Token invalid")
		return claims, constant.ErrUnAuthentication
	}

	return claims, nil
}
