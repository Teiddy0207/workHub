package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"

	"workHub/constant"
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
			fmt.Println("JWT validation error:", ve)
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				fmt.Println("Token expired")
				return claims, constant.ErrUnAuthentication
			}
		} else {
			fmt.Println("JWT parse error:", err)
		}
		return claims, constant.ErrUnAuthentication
	}

	if !token.Valid {
		fmt.Println("Token invalid")
		return claims, constant.ErrUnAuthentication
	}

	return claims, nil
}
