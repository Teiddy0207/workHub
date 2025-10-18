package utils

import (
	"fmt"
	"strings"

	"workHub/constant"
	"workHub/internal/entity"

	"gorm.io/gorm"
)

func ApplyFilters(db *gorm.DB, filters []entity.FindRequestFilter) *gorm.DB {
	for _, filter := range filters {
		db = applySingleFilter(db, filter)
	}
	return db
}

func applySingleFilter(db *gorm.DB, filter entity.FindRequestFilter) *gorm.DB {
	if filter.IgnorePrepare {
		return db
	}

	// Custom function hỗ trợ linh hoạt
	if filter.CustomFunc != nil {
		sql, args := filter.CustomFunc()
		if argsSlice, ok := args.([]interface{}); ok {
			return db.Where(sql, argsSlice...)
		}
		return db.Where(sql, args)
	}

	// SubFilters: mặc định nối bằng AND
	if len(filter.SubFilters) > 0 {
		return db.Where(func(tx *gorm.DB) *gorm.DB {
			for _, sub := range filter.SubFilters {
				tx = applySingleFilter(tx, sub)
			}
			return tx
		})
	}

	// Default operator
	key := filter.Key
	val := filter.Value

	// Validate column name to prevent SQL injection
	sanitizedKey, err := SanitizeColumnName(key)
	if err != nil {
		// Return error or skip this filter
		return db
	}

	switch strings.ToLower(filter.Operator) {
	case constant.OPERATOR_EQUAL, "":
		return db.Where(fmt.Sprintf("%s = ?", sanitizedKey), val)
	case constant.OPERATOR_NOT_EQUAL:
		return db.Where(fmt.Sprintf("%s <> ?", sanitizedKey), val)
	case constant.OPERATOR_IN:
		return db.Where(fmt.Sprintf("%s IN ?", sanitizedKey), val)
	case constant.OPERATOR_NOT_IN:
		return db.Where(fmt.Sprintf("%s NOT IN ?", sanitizedKey), val)
	case constant.OPERATOR_LIKE:
		return db.Where(fmt.Sprintf("%s LIKE ?", sanitizedKey), val)
	case constant.OPERATOR_ILIKE:
		return db.Where(fmt.Sprintf("%s ILIKE ?", sanitizedKey), val) // Postgres only
	case constant.OPERATOR_GREATER:
		return db.Where(fmt.Sprintf("%s > ?", sanitizedKey), val)
	case constant.OPERATOR_GREATER_EQ:
		return db.Where(fmt.Sprintf("%s >= ?", sanitizedKey), val)
	case constant.OPERATOR_LESS:
		return db.Where(fmt.Sprintf("%s < ?", sanitizedKey), val)
	case constant.OPERATOR_LESS_EQ:
		return db.Where(fmt.Sprintf("%s <= ?", sanitizedKey), val)
	default:
		return db.Where(fmt.Sprintf("%s = ?", sanitizedKey), val)
	}

}
