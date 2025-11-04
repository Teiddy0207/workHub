package router

import (
	"fmt"
	internalconfig "workHub/internal/config"
	"workHub/internal/service"
	"workHub/logger"
	"workHub/pkg/jwt"

	"crypto/rsa"
)


type JWTConfig struct {
	Service   service.JWTServiceInterface
	PublicKey *rsa.PublicKey
}


func InitJWTConfig() (*JWTConfig, error) {
	cfg, err := internalconfig.LoadConfig()
	if err != nil {
		logger.Error("config", "LoadConfig", fmt.Sprintf("Config load failed: %v", err))
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	jwtService, err := service.NewJWTService(*cfg)
	if err != nil {
		logger.Error("service", "NewJWTService", fmt.Sprintf("JWT service init failed: %v", err))
		return nil, fmt.Errorf("failed to init JWT service: %w", err)
	}

	_, publicKey, _, err := jwt.ParseKey(cfg.Jwt)
	if err != nil {
		logger.Error("router", "InitJWTConfig", fmt.Sprintf("Failed to parse JWT keys: %v", err))
		return nil, fmt.Errorf("failed to parse JWT keys: %w", err)
	}

	return &JWTConfig{
		Service:   jwtService,
		PublicKey: publicKey,
	}, nil
}

