package service

import (
	"context"
	"fmt"
	"time"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/pkg/utils"
	"workHub/constant"
)

type AuthService struct {
	AuthRepo   repository.AuthRepository
	
}

type AuthServiceInterface interface {
	GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

func NewAuthService(AuthRepo repository.AuthRepository) AuthServiceInterface {
	return &AuthService{
		AuthRepo:   AuthRepo,
	}
}

func (service *AuthService) GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error) {
	users, err := service.AuthRepo.ListUsers(ctx, params)

	if err != nil {
		return dto.PaginatedUserResponse{}, err
	}

	return mapper.ToRegisterResponse(users), nil

}

func (service *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	fmt.Printf("🔍 Login attempt for email: %s\n", req.Email)
	
	// 1. Tìm user theo email
	user, err := service.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		fmt.Printf("❌ User not found: %v\n", err)
		return dto.LoginResponse{}, constant.ErrUnAuthentication
	}
	
	fmt.Printf("✅ User found: %s (ID: %s)\n", user.Username, user.ID)

	// 2. Verify password
	err = utils.CompareHashPassword(req.Password, user.Password)
	if err != nil {
		fmt.Printf("❌ Password incorrect: %v\n", err)
		return dto.LoginResponse{}, constant.ErrUnAuthentication
	}
	
	fmt.Printf("✅ Password verified successfully\n")

	// 3. Tạo response đơn giản (tạm thời không dùng JWT vì cần config)
	// TODO: Implement JWT token generation sau khi có RSA keys
	
	// 4. Tạo response
	response := dto.LoginResponse{
		AccessToken:  "temp_access_token_" + user.ID,
		RefreshToken: "temp_refresh_token_" + user.ID,
		ExpiresAt:    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		User: dto.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	
	fmt.Printf("🎉 Login successful, returning response\n")
	return response, nil
}
