package service

import (
	"context"
	"time"
	"workHub/internal/config"
	"workHub/internal/dto"
	"workHub/pkg/jwt"
	"crypto/rsa"
	jwtgo "github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	signMethod jwtgo.SigningMethod
	config     jwt.Config
}

type JWTServiceInterface interface {
	GenerateAccessToken(ctx context.Context, userInfo dto.Users) (string, *time.Time, error)
	GenerateRefreshToken(ctx context.Context, userInfo dto.Users) (string, *time.Time, error)
}

func NewJWTService(cfg config.Config) (JWTServiceInterface, error) {
	privateKey, publicKey, signMethod, err := jwt.ParseKey(cfg.Jwt)
	if err != nil {
		return nil, err
	}

	return &JWTService{
		privateKey: privateKey,
		publicKey:  publicKey,
		signMethod: signMethod,
		config:     cfg.Jwt,
	}, nil
}

func (j *JWTService) GenerateAccessToken(ctx context.Context, userInfo dto.Users) (string, *time.Time, error) {
	jwtReq := jwt.JwtReq{
		UserInfo: userInfo,
	}

	return jwt.GenerateToken(
		ctx,
		jwtReq,
		j.signMethod,
		j.privateKey,
		jwt.TokenTypeAccessToken,
		j.config.Issuer,
		j.config.ShortTokenExpireTime,
	)
}

func (j *JWTService) GenerateRefreshToken(ctx context.Context, userInfo dto.Users) (string, *time.Time, error) {
	jwtReq := jwt.JwtReq{
		UserInfo: userInfo,
	}

	return jwt.GenerateToken(
		ctx,
		jwtReq,
		j.signMethod,
		j.privateKey,
		jwt.TokenTypeRefreshToken,
		j.config.Issuer,
		j.config.RefreshTokenExpireTime,
	)
}
