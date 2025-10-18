package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(
	ctx context.Context,
	req JwtReq,
	signMethod jwt.SigningMethod,
	privateKey *rsa.PrivateKey,
	tokenType TokenType,
	Issuer string,
	expiredTime uint,
) (string, *time.Time, error) {
	var (
		tokenStr string
		err      error
	)
	expiredAt := time.Now().UTC().Add(time.Second * time.Duration(expiredTime))
	jwtToken := jwt.New(signMethod)

	jwtClaim := JwtClaim{
		UserInfo: req.UserInfo,
		Type:     tokenType,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			Issuer:    Issuer,
		},
	}
	jwtToken.Claims = jwtClaim
	tokenStr, err = jwtToken.SignedString(privateKey)
	if err != nil {
		return tokenStr, nil, err
	}
	return tokenStr, &expiredAt, nil
}

func ParseToken(
	ctx context.Context,
	tokenStr string,
	publicKey *rsa.PublicKey,
) (*JwtClaim, error) {
	// parse và verify token
	parsedToken, err := jwt.ParseWithClaims(tokenStr, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		// đảm bảo đúng thuật toán
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// lấy claims ra
	if claims, ok := parsedToken.Claims.(*JwtClaim); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func ParseRSAPublicKeyFromPEM(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return rsaPub, nil
}
