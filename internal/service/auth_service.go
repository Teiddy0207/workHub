package service

import (
	"context"
	"fmt"
	"time"
	"workHub/internal/dto"
	"workHub/internal/entity"
	"workHub/internal/mapper"
	"workHub/internal/repository"
	"workHub/pkg/params"
	"workHub/pkg/utils"
	"workHub/constant"
	"workHub/logger"
	"github.com/google/uuid"
)

type AuthService struct {
	AuthRepo          repository.AuthRepository
	SessionRepo       repository.SessionRepository
	SessionRedisRepo  repository.SessionRedisRepository
	JWTService        JWTServiceInterface
}

type AuthServiceInterface interface {
	GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, ipAddress, userAgent string) (dto.LoginResponse, error)
}

func NewAuthService(AuthRepo repository.AuthRepository, SessionRepo repository.SessionRepository, SessionRedisRepo repository.SessionRedisRepository, JWTService JWTServiceInterface) AuthServiceInterface {
	return &AuthService{
		AuthRepo:         AuthRepo,
		SessionRepo:      SessionRepo,
		SessionRedisRepo: SessionRedisRepo,
		JWTService:       JWTService,
	}
}

func (service *AuthService) GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error) {
	users, err := service.AuthRepo.ListUsers(ctx, params)

	if err != nil {
		return dto.PaginatedUserResponse{}, err
	}

	return mapper.ToRegisterResponse(users), nil

}

func (service *AuthService) Login(ctx context.Context, req dto.LoginRequest, ipAddress, userAgent string) (dto.LoginResponse, error) {
	logger.Info("service", "Login", fmt.Sprintf("Login attempt for email: %s", req.Email))
	
	user, err := service.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("service", "Login", fmt.Sprintf("User not found: %v", err))
		return dto.LoginResponse{}, constant.ErrUsernameOrPasswordIncorrect
	}
	
	logger.Info("service", "Login", fmt.Sprintf("User found: %s (ID: %s)", user.Username, user.ID))

	err = utils.CompareHashPassword(req.Password, user.Password)
	if err != nil {
		logger.Error("service", "Login", fmt.Sprintf("Password incorrect: %v", err))
		return dto.LoginResponse{}, constant.ErrPasswordIncorrect
	}
	
	logger.Info("service", "Login", "Password verified successfully")

	// Tạo JWT tokens
	accessToken, accessExpiresAt, err := service.JWTService.GenerateAccessTokenFromEntity(ctx, user)
	if err != nil {
		logger.Error("service", "Login", fmt.Sprintf("Failed to generate access token: %v", err))
		return dto.LoginResponse{}, constant.ErrInternalServer
	}

	refreshToken, _, err := service.JWTService.GenerateRefreshTokenFromEntity(ctx, user)
	if err != nil {
		logger.Error("service", "Login", fmt.Sprintf("Failed to generate refresh token: %v", err))
		return dto.LoginResponse{}, constant.ErrInternalServer
	}

	// Lưu session vào database
	session := &entity.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    *accessExpiresAt, // Dereference pointer
		IsActive:     true,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}
	
	// Lưu session vào database
	err = service.SessionRepo.CreateSession(ctx, session)
	if err != nil {
		logger.Error("service", "Login", fmt.Sprintf("Failed to create session: %v", err))
		// Không return error vì tokens đã được tạo, chỉ log warning
		logger.Warn("service", "Login", "Session not saved to DB but login continues")
	} else {
		logger.Info("service", "Login", fmt.Sprintf("Session saved to DB successfully: %s", session.ID))
	}

	// Lưu session vào Redis cache
	if service.SessionRedisRepo != nil {
		expiration := time.Until(*accessExpiresAt)
		if expiration > 0 {
			err = service.SessionRedisRepo.SaveSession(ctx, session, expiration)
			if err != nil {
				logger.Error("service", "Login", fmt.Sprintf("Failed to save session to Redis: %v", err))
				// Không return error, chỉ log warning vì session đã được lưu vào DB
				logger.Warn("service", "Login", "Session not saved to Redis but login continues")
			} else {
				logger.Info("service", "Login", fmt.Sprintf("Session saved to Redis successfully: %s", session.ID))
			}
		}
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
	
	logger.Info("service", "Login", "Login successful, JWT tokens generated")
	return response, nil
}
