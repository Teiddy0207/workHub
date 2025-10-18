package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidColumnNames contains all valid column names that can be used in queries
var ValidColumnNames = map[string]bool{
	// User table columns
	"id": true, "name": true, "username": true, "email": true, "phone": true,
	"password": true, "status": true, "role": true, "role_id": true,
	"department": true, "department_id": true, "store_name": true, "location_id": true,
	"updated_by": true, "avatar_url": true, "avatar": true, "gender": true,
	"dob": true, "bio": true, "is_vip": true, "created_at": true,
	"updated_at": true, "deleted_at": true,

	// Queue ticket table columns
	"store_id": true, "user_id": true, "queue_id": true,
	"ticket_number": true, "customer_name": true, "customer_phone": true,

	// Store table columns
	"telecom_service_point_name": true, "telecom_service_point_phone": true,
	"telecom_service_point_address": true, "telecom_service_point_name_unaccent": true,

	// Role table columns
	"code": true, "description": true, "is_active": true,

	// Permission table columns
	"module": true,

	// Service types table columns
	"name_unaccent": true,

	// Feedback table columns
	"rating": true, "notes": true, "queue_ticket_id": true,
}

// ValidOrderDirections contains valid ORDER BY directions
var ValidOrderDirections = map[string]bool{
	"ASC":  true,
	"DESC": true,
}

// IsValidColumnName checks if a column name is safe to use in SQL queries
func IsValidColumnName(columnName string) bool {
	if columnName == "" {
		return false
	}
	
	// Check if it's in our whitelist
	if ValidColumnNames[strings.ToLower(columnName)] {
		return true
	}
	
	// Additional validation: only allow alphanumeric characters and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, columnName)
	return matched
}

// IsValidOrderDirection checks if an order direction is valid
func IsValidOrderDirection(direction string) bool {
	return ValidOrderDirections[strings.ToUpper(direction)]
}

// SanitizeColumnName validates and returns a safe column name
func SanitizeColumnName(columnName string) (string, error) {
	if !IsValidColumnName(columnName) {
		return "", fmt.Errorf("invalid column name: %s", columnName)
	}
	return strings.ToLower(columnName), nil
}

// SanitizeOrderDirection validates and returns a safe order direction
func SanitizeOrderDirection(direction string) (string, error) {
	if !IsValidOrderDirection(direction) {
		return "", fmt.Errorf("invalid order direction: %s", direction)
	}
	return strings.ToUpper(direction), nil
}

// ValidateFieldName checks if a field name is safe for credential checking
func ValidateFieldName(field string) error {
	validFields := map[string]bool{
		"email":    true,
		"phone":    true,
		"username": true,
	}
	
	if !validFields[strings.ToLower(field)] {
		return fmt.Errorf("invalid field name: %s", field)
	}
	return nil
}

// ValidateUploadFile validates uploaded file for security and format
func ValidateUploadFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return fmt.Errorf("file header is nil")
	}

	// Check file size (max 10MB)
	const maxSize = 10 * 1024 * 1024
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	// Check file extension
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".pdf":  true,
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	// Check MIME type
	allowedMimeTypes := map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
		"application/pdf": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		return fmt.Errorf("unsupported MIME type: %s", contentType)
	}

	return nil
} 