package repository

import (
	"context"
	"fmt"
	"workHub/internal/entity"
	"workHub/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSessionByToken(ctx context.Context, accessToken string) (entity.Session, error)
	GetActiveSessionsByUserID(ctx context.Context, userID string) ([]entity.Session, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	logger.Info("repository", "CreateSession", fmt.Sprintf("Creating session for user: %s", session.UserID))
	
	// Generate UUID nếu chưa có
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	
	err := r.db.WithContext(ctx).Create(session).Error
	if err != nil {
		logger.Error("repository", "CreateSession", fmt.Sprintf("Failed to create session: %v", err))
		return err
	}
	
	logger.Info("repository", "CreateSession", fmt.Sprintf("Session created successfully: %s", session.ID))
	return nil
}

func (r *sessionRepository) GetSessionByToken(ctx context.Context, accessToken string) (entity.Session, error) {
	var session entity.Session
	err := r.db.WithContext(ctx).
		Where("access_token = ? AND is_active = ?", accessToken, true).
		First(&session).Error
	
	if err != nil {
		return entity.Session{}, err
	}
	
	return session, nil
}

func (r *sessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID string) ([]entity.Session, error) {
	var sessions []entity.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Find(&sessions).Error
	
	return sessions, err
}

func (r *sessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	logger.Info("repository", "RevokeSession", fmt.Sprintf("Revoking session: %s", sessionID))
	
	result := r.db.WithContext(ctx).
		Model(&entity.Session{}).
		Where("id = ?", sessionID).
		Update("is_active", false)
	
	if result.Error != nil {
		logger.Error("repository", "RevokeSession", fmt.Sprintf("Failed to revoke session: %v", result.Error))
		return result.Error
	}
	
	logger.Info("repository", "RevokeSession", fmt.Sprintf("Session revoked: %s", sessionID))
	return nil
}

func (r *sessionRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	logger.Info("repository", "RevokeAllUserSessions", fmt.Sprintf("Revoking all sessions for user: %s", userID))
	
	result := r.db.WithContext(ctx).
		Model(&entity.Session{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false)
	
	if result.Error != nil {
		logger.Error("repository", "RevokeAllUserSessions", fmt.Sprintf("Failed to revoke sessions: %v", result.Error))
		return result.Error
	}
	
	logger.Info("repository", "RevokeAllUserSessions", fmt.Sprintf("All sessions revoked for user: %s", userID))
	return nil
}

