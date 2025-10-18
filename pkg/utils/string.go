package utils

import (
	"fmt"
	"mime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func ToUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func RemoveAccentVietnamese(s string) string {
	s = strings.ReplaceAll(s, "Đ", "D") // Don't know why Đ|đ cannot be transformed to D|d so we need to use this
	s = strings.ReplaceAll(s, "đ", "d")
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	res, _, _ := transform.String(t, s)
	return res
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty checks if a string is empty after trimming spaces
func IsEmpty(s string) bool {
	return TrimSpace(s) == ""
}

// GenerateFileName generates a unique filename with timestamp and UUID
func GenerateFileName(originalName, ext string) string {
	timestamp := time.Now().Unix()
	uuid := uuid.New().String()
	return fmt.Sprintf("%d_%s%s", timestamp, uuid, ext)
}

// GetExtensionFromContentType gets file extension from MIME content type
func GetExtensionFromContentType(contentType string) string {
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil || len(ext) == 0 {
		return ""
	}
	return ext[0]
}
