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
	JWTService JWTServiceInterface
}

type AuthServiceInterface interface {
	GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

func NewAuthService(AuthRepo repository.AuthRepository, JWTService JWTServiceInterface) AuthServiceInterface {
	return &AuthService{
		AuthRepo:   AuthRepo,
		JWTService: JWTService,
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
	
	user, err := service.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		fmt.Printf("❌ User not found: %v\n", err)
		return dto.LoginResponse{}, constant.ErrUsernameOrPasswordIncorrect
	}
	
	fmt.Printf("✅ User found: %s (ID: %s)\n", user.Username, user.ID)

	err = utils.CompareHashPassword(req.Password, user.Password)
	if err != nil {
		fmt.Printf("❌ Password incorrect: %v\n", err)
		return dto.LoginResponse{}, constant.ErrPasswordIncorrect
	}
	
	fmt.Printf("✅ Password verified successfully\n")

	// Tạo JWT tokens thực sự
	accessToken, accessExpiresAt, err := service.JWTService.GenerateAccessTokenFromEntity(ctx, user)
	if err != nil {
		fmt.Printf("❌ Failed to generate access token: %v\n", err)
		return dto.LoginResponse{}, constant.ErrInternalServer
	}

	refreshToken, _, err := service.JWTService.GenerateRefreshTokenFromEntity(ctx, user)
	if err != nil {
		fmt.Printf("❌ Failed to generate refresh token: %v\n", err)
		return dto.LoginResponse{}, constant.ErrInternalServer
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Format(time.RFC3339),
		User: dto.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	
	fmt.Printf("✅ Login successful, JWT tokens generated\n")
	return response, nil
}
