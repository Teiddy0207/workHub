package jwt

import (
	"workHub/internal/dto"
	"github.com/dgrijalva/jwt-go"
)

type TokenType int8

const (
	TokenTypeAccessToken  TokenType = 1
	TokenTypeRefreshToken TokenType = 2
)

type JwtReq struct {
	UserInfo dto.Users
}

type JwtClaim struct {
	*jwt.StandardClaims
	UserInfo dto.Users
	Type TokenType
}

type JwtConfig struct {
	jwt.StandardClaims
	SigningMethod string
	PublicKey     string
	PrivateKey    string
}

type Config struct {
	SigningMethod          string `yaml:"signing_method" mapstructure:"signing_method"`
	PrivateKey             string `yaml:"private_key" mapstructure:"private_key"`
	PublicKey              string `yaml:"public_key" mapstructure:"public_key"`
	Issuer                 string `yaml:"issuer" mapstructure:"issuer"`
	RefreshTokenExpireTime uint   `yaml:"refresh_token_expire" mapstructure:"refresh_token_expire"`
	LongTokenExpireTime    uint   `yaml:"long_token_expire" mapstructure:"long_token_expire"`
	ShortTokenExpireTime   uint   `yaml:"short_token_expire" mapstructure:"short_token_expire"`
	IsRefreshToken         bool   `yaml:"is_refresh_token" mapstructure:"is_refresh_token"`
	ValidatePassword       bool   `yaml:"validate_password" mapstructure:"validate_password"`
	LenToken               int    `yaml:"len_token" mapstructure:"len_token"`
	TokenExpire            int    `yaml:"token_expire" mapstructure:"token_expire"`
}
