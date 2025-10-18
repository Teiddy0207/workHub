package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func ParseKey(cfg Config) (*rsa.PrivateKey, *rsa.PublicKey, jwt.SigningMethod, error) {
    var (
        private *rsa.PrivateKey
        public  *rsa.PublicKey
        sign    jwt.SigningMethod
        err     error
    )

    privateByte, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to decode private key: %w", err)
    }

    private, err = jwt.ParseRSAPrivateKeyFromPEM(privateByte)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to parse private key from PEM: %w", err)
    }

    publicByte, err := base64.StdEncoding.DecodeString(cfg.PublicKey)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to decode public key: %w", err)
    }

    public, err = jwt.ParseRSAPublicKeyFromPEM(publicByte)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to parse public key from PEM: %w", err)
    }

    sign = jwt.GetSigningMethod(cfg.SigningMethod)
    if sign == nil {
        return nil, nil, nil, fmt.Errorf("invalid signing method: %s", cfg.SigningMethod)
    }

    return private, public, sign, nil
}