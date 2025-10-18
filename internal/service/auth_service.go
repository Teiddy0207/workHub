package service

import (
	"context"
	"time"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/pkg/jwt"
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
	// 1. Tìm user theo email
	user, err := service.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, constant.ErrUnAuthentication
	}

	// 2. Verify password
	err = utils.CompareHashPassword(req.Password, user.Password)
	if err != nil {
		return dto.LoginResponse{}, constant.ErrUnAuthentication
	}

	// 3. Tạo JWT tokens
	userInfo := dto.Users{
		Username: user.Username,
		Email:    user.Email,
	}

	jwtReq := jwt.JwtReq{
		UserInfo: userInfo,
	}

	// TODO: Cần load config từ file hoặc environment variables
	// Tạm thời sử dụng hardcode values
	accessToken, expiresAt, err := jwt.GenerateToken(
		ctx,
		jwtReq,
		nil, // TODO: Get signing method from config
		nil, // TODO: Get private key from config
		jwt.TokenTypeAccessToken,
		"workHub",
		3600, // 1 hour
	)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, _, err := jwt.GenerateToken(
		ctx,
		jwtReq,
		nil, // TODO: Get signing method from config
		nil, // TODO: Get private key from config
		jwt.TokenTypeRefreshToken,
		"workHub",
		604800, // 7 days
	)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// 4. Tạo response
	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		User: dto.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}, nil
}
