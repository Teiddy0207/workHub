package mapper

import (
	"time"

	"workHub/internal/dto"
	"workHub/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ToRegisterResponse(user entity.PaginatedUsers) dto.PaginatedUserResponse {
	var response []dto.RegisterResponse

	for _, u := range user.Items {
		response = append(response, dto.RegisterResponse{
			ID:          u.ID,
			Email:       u.Email,
			FullName:    u.FullName,
			PhoneNumber: u.PhoneNumber,
			Address:     u.Address,
			Role:        u.Role,
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		})
	}
	var totalPages int
	if user.PageSize > 0 {
		totalPages = user.TotalItems / user.PageSize
		if user.TotalItems%user.PageSize > 0 {
			totalPages++
		}
	} else {
		totalPages = 1
	}
	return dto.PaginatedUserResponse{
		Items:       response,
		TotalItems:  user.TotalItems,
		TotalPages:  totalPages,
		CurrentPage: user.PageNumber,
		PageSize:    user.PageSize,
	}
}

func ToUserEntity(req dto.RegisterRequest) (entity.User, error) {
	hashed, err := HashPassword(req.Password)
	if err != nil {
		return entity.User{}, err
	}
	role := req.Role
	if role == "" {
		role = "student"
	}
	return entity.User{
		ID:          uuid.NewString(),
		Email:       req.Email,
		Username:    req.Email, // Use email as username if not provided
		Password:    hashed,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Role:        role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func ToUserResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}
}

func ToUserInfo(user entity.User) dto.UserInfo {
	return dto.UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		Role:        user.Role,
	}
}

func ToUserItem(u *entity.User) dto.UserItem {
	return dto.UserItem{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
