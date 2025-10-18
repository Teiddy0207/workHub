package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(val string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.Trim(val, " ")), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashPassword(val string, hashVal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashVal), []byte(strings.Trim(val, " ")))
}

func EncryptString(text, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(cryptoText, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// HashRandomStringBase64 generates a base64 hash of a random string
func HashRandomStringBase64(input string) string {
	// Convert input to bytes and encode to base64
	hash := base64.URLEncoding.EncodeToString([]byte(input))
	return hash
}

func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateRandomStringHashOnly generates a random string and returns only alphanumeric characters
func GenerateRandomStringHashOnly(length int) (string, error) {
	// Generate random string using crypto/rand with only alphanumeric characters
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}

	// Return the random string directly (no hashing needed)
	return string(bytes), nil
}

// generateChannelHash generates a random hash string with ws_channel prefix
func GenerateChannelHash() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return "ws_channel" + hex.EncodeToString(bytes)
}