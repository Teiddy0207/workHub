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

	for _, users := range user.Items {
		response = append(response, dto.RegisterResponse{
			ID:       users.ID,
			Email:    users.Email,
			Username: users.Username,
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
	return entity.User{
		ID:        uuid.NewString(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashed,
		CreatedAt: time.Now(),
	}, nil
}

func ToUserItem(u *entity.User) dto.UserItem {
	return dto.UserItem{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
