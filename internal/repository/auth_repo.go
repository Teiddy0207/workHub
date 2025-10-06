package repository

import (
	"workHub/internal/entity"
	"workHub/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepository interface {
    CreateUser(user *entity.User) (*entity.User, error)
    GetUserByEmail(email string) (*entity.User, error)
    GetUserByUsername(username string) (*entity.User, error)
    ListUsers(keyword string, page, limit int) ([]entity.User, int64, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		logger.Log.Error("auth_repo:CreateUser", zap.Error(err))
		return nil, err
	}
	
	logger.Log.Info("auth_repo:CreateUser", zap.String("email", user.Email))
	return user, nil
}

func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) ListUsers(keyword string, page, limit int) ([]entity.User, int64, error) {
    var users []entity.User
    var total int64

    q := r.db.Model(&entity.User{})
    if keyword != "" {
        like := "%" + keyword + "%"
        q = q.Where("email ILIKE ? OR username ILIKE ?", like, like)
    }

    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit
    if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
        return nil, 0, err
    }
    return users, total, nil
}