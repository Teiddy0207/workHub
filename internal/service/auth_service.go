package service

import (
	"context"
	"workHub/internal/dto"
	"workHub/internal/mapper"
	"workHub/internal/repository"

	// "workHub/internal/service"
	"workHub/pkg/params"
)

type AuthService struct {
	AuthRepo repository.AuthRepository
}

type AuthServiceInterface interface {
	GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error)
}

func NewAuthService(AuthRepo repository.AuthRepository) AuthServiceInterface {
	return &AuthService{
		AuthRepo: AuthRepo,
	}
}

func (service *AuthService) GetListUser(ctx context.Context, params params.QueryParams) (dto.PaginatedUserResponse, error) {
	users, err := service.AuthRepo.ListUsers(ctx, params)

	if err != nil {
		return dto.PaginatedUserResponse{}, err
	}

	return mapper.ToRegisterResponse(users), nil

}
