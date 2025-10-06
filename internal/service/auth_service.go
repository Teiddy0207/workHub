package service

import (
    "context"
    "errors"
    "go.uber.org/zap"
    "workHub/internal/dto"
    "workHub/internal/mapper"
    "workHub/internal/repository"
    "workHub/logger"
    "workHub/utils"
)

type AuthService interface {
    Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
    GetListUser(ctx context.Context, keyword string, page, limit int) ([]dto.UserItem, utils.Pagination, error)
}

type DefaultAuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *DefaultAuthService {
	return &DefaultAuthService{repo: repo}
}

func (s *DefaultAuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	log := logger.WithTrace(ctx, "auth_service", "Register")
	log.Info("start register", zap.String("email", req.Email))

	if err := s.checkDuplicate(ctx, req); err != nil {
		return dto.RegisterResponse{}, err
	}

	user, err := mapper.ToUserEntity(req)
	if err != nil {
		log.Error("mapper failed", zap.Error(err))
		return dto.RegisterResponse{}, err
	}

	created, err := s.repo.CreateUser(&user)
	if err != nil {
		log.Error("repo create failed", zap.Error(err))
		return dto.RegisterResponse{}, err
	}

	log.Info("register success", zap.String("email", created.Email))
	return mapper.ToRegisterResponse(created), nil
}




func (s *DefaultAuthService) checkDuplicate(ctx context.Context, req dto.RegisterRequest) error {
	log := logger.WithTrace(ctx, "auth_service", "checkDuplicate")

	if u, err := s.repo.GetUserByEmail(req.Email); err != nil {
		log.Error("query email failed", zap.Error(err))
		return err
	} else if u != nil {
		log.Warn("duplicate email", zap.String("email", req.Email))
		return errors.New("email already exists")
	}

	if u, err := s.repo.GetUserByUsername(req.Username); err != nil {
		log.Error("query username failed", zap.Error(err))
		return err
	} else if u != nil {
		log.Warn("duplicate username", zap.String("username", req.Username))
		return errors.New("username already exists")
	}

	return nil
}

func (s *DefaultAuthService) GetListUser(ctx context.Context, keyword string, page, limit int) ([]dto.UserItem, utils.Pagination, error) {
    log := logger.WithTrace(ctx, "auth_service", "GetListUser")
    log.Info("start list users", zap.String("keyword", keyword), zap.Int("page", page), zap.Int("limit", limit))

    users, total, err := s.repo.ListUsers(keyword, page, limit)
    if err != nil {
        log.Error("repo list failed", zap.Error(err))
        return nil, utils.Paginate(page, limit, 0), err
    }

    items := make([]dto.UserItem, 0, len(users))
    for i := range users {
        items = append(items, mapper.ToUserItem(&users[i]))
    }

    meta := utils.Paginate(page, limit, int(total))
    log.Info("list users success", zap.Int("count", len(items)), zap.Int("total", int(total)))
    
	
	return items, meta, nil
}