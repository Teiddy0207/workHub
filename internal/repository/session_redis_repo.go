package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"workHub/internal/entity"
	"workHub/pkg/utils"
	"workHub/logger"
)

type SessionRedisRepository interface {
	SaveSession(ctx context.Context, session *entity.Session, expiration time.Duration) error
	GetSessionByToken(ctx context.Context, accessToken string) (*entity.Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteSessionByToken(ctx context.Context, accessToken string) error
	DeleteUserSessions(ctx context.Context, userID string) error
	IsSessionActive(ctx context.Context, accessToken string) (bool, error)
}

type sessionRedisRepository struct{}

func NewSessionRedisRepository() SessionRedisRepository {
	return &sessionRedisRepository{}
}

// SaveSession lưu session vào Redis với key pattern: session:{sessionID}
// Và thêm mapping token -> sessionID: token:{accessToken} -> sessionID
func (r *sessionRedisRepository) SaveSession(ctx context.Context, session *entity.Session, expiration time.Duration) error {
	if utils.RedisClient == nil {
		logger.Warn("repository", "SaveSession", "Redis client is nil, skipping Redis save")
		return nil
	}

	logger.Info("repository", "SaveSession", fmt.Sprintf("Saving session to Redis: %s", session.ID))

	// Serialize session to JSON
	sessionData, err := json.Marshal(session)
	if err != nil {
		logger.Error("repository", "SaveSession", fmt.Sprintf("Failed to marshal session: %v", err))
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Lưu session với key: session:{sessionID}
	sessionKey := fmt.Sprintf("session:%s", session.ID)
	err = utils.SaveToRedis(ctx, utils.RedisClient, sessionKey, string(sessionData), expiration)
	if err != nil {
		logger.Error("repository", "SaveSession", fmt.Sprintf("Failed to save session to Redis: %v", err))
		return err
	}

	// Lưu mapping token -> sessionID với key: token:{accessToken}
	tokenKey := fmt.Sprintf("token:%s", session.AccessToken)
	err = utils.SaveToRedis(ctx, utils.RedisClient, tokenKey, session.ID, expiration)
	if err != nil {
		logger.Error("repository", "SaveSession", fmt.Sprintf("Failed to save token mapping: %v", err))
		return err
	}

	// Lưu danh sách session IDs của user với key: user:{userID}:sessions
	userSessionsKey := fmt.Sprintf("user:%s:sessions", session.UserID)
	// Thêm session ID vào set (sử dụng SAdd)
	err = utils.RedisClient.SAdd(ctx, userSessionsKey, session.ID).Err()
	if err != nil {
		logger.Error("repository", "SaveSession", fmt.Sprintf("Failed to add session to user set: %v", err))
		return err
	}
	// Set expiration cho user sessions set
	err = utils.RedisClient.Expire(ctx, userSessionsKey, expiration).Err()
	if err != nil {
		logger.Warn("repository", "SaveSession", fmt.Sprintf("Failed to set expiration for user sessions: %v", err))
	}

	logger.Info("repository", "SaveSession", fmt.Sprintf("Session saved to Redis successfully: %s", session.ID))
	return nil
}

// GetSessionByToken lấy session từ Redis bằng access token
func (r *sessionRedisRepository) GetSessionByToken(ctx context.Context, accessToken string) (*entity.Session, error) {
	if utils.RedisClient == nil {
		return nil, fmt.Errorf("Redis client is nil")
	}

	// Lấy sessionID từ token mapping
	tokenKey := fmt.Sprintf("token:%s", accessToken)
	sessionID, err := utils.GetFromRedis(ctx, utils.RedisClient, tokenKey)
	if err != nil {
		logger.Warn("repository", "GetSessionByToken", fmt.Sprintf("Token not found in Redis: %v", err))
		return nil, fmt.Errorf("token not found: %w", err)
	}

	// Lấy session từ sessionID
	return r.GetSessionByID(ctx, sessionID)
}

// GetSessionByID lấy session từ Redis bằng session ID
func (r *sessionRedisRepository) GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error) {
	if utils.RedisClient == nil {
		return nil, fmt.Errorf("Redis client is nil")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	sessionData, err := utils.GetFromRedis(ctx, utils.RedisClient, sessionKey)
	if err != nil {
		logger.Warn("repository", "GetSessionByID", fmt.Sprintf("Session not found in Redis: %v", err))
		return nil, fmt.Errorf("session not found: %w", err)
	}

	var session entity.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		logger.Error("repository", "GetSessionByID", fmt.Sprintf("Failed to unmarshal session: %v", err))
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// DeleteSession xóa session khỏi Redis
func (r *sessionRedisRepository) DeleteSession(ctx context.Context, sessionID string) error {
	if utils.RedisClient == nil {
		return nil
	}

	logger.Info("repository", "DeleteSession", fmt.Sprintf("Deleting session from Redis: %s", sessionID))

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	
	// Lấy session để xóa token mapping và user sessions
	session, err := r.GetSessionByID(ctx, sessionID)
	if err == nil && session != nil {
		// Xóa token mapping
		tokenKey := fmt.Sprintf("token:%s", session.AccessToken)
		err = utils.RedisClient.Del(ctx, tokenKey).Err()
		if err != nil {
			logger.Warn("repository", "DeleteSession", fmt.Sprintf("Failed to delete token mapping: %v", err))
		}

		// Xóa khỏi user sessions set
		userSessionsKey := fmt.Sprintf("user:%s:sessions", session.UserID)
		err = utils.RedisClient.SRem(ctx, userSessionsKey, sessionID).Err()
		if err != nil {
			logger.Warn("repository", "DeleteSession", fmt.Sprintf("Failed to remove from user sessions: %v", err))
		}
	}

	// Xóa session
	err = utils.RedisClient.Del(ctx, sessionKey).Err()
	if err != nil {
		logger.Error("repository", "DeleteSession", fmt.Sprintf("Failed to delete session: %v", err))
		return err
	}

	logger.Info("repository", "DeleteSession", fmt.Sprintf("Session deleted from Redis: %s", sessionID))
	return nil
}

// DeleteSessionByToken xóa session khỏi Redis bằng access token
func (r *sessionRedisRepository) DeleteSessionByToken(ctx context.Context, accessToken string) error {
	if utils.RedisClient == nil {
		return nil
	}

	// Lấy sessionID từ token
	session, err := r.GetSessionByToken(ctx, accessToken)
	if err != nil {
		return err
	}

	return r.DeleteSession(ctx, session.ID)
}

// DeleteUserSessions xóa tất cả sessions của user
func (r *sessionRedisRepository) DeleteUserSessions(ctx context.Context, userID string) error {
	if utils.RedisClient == nil {
		return nil
	}

	logger.Info("repository", "DeleteUserSessions", fmt.Sprintf("Deleting all sessions for user: %s", userID))

	userSessionsKey := fmt.Sprintf("user:%s:sessions", userID)
	
	// Lấy tất cả session IDs từ set
	sessionIDs, err := utils.RedisClient.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		logger.Error("repository", "DeleteUserSessions", fmt.Sprintf("Failed to get user sessions: %v", err))
		return err
	}

	// Xóa từng session
	for _, sessionID := range sessionIDs {
		if err := r.DeleteSession(ctx, sessionID); err != nil {
			logger.Warn("repository", "DeleteUserSessions", fmt.Sprintf("Failed to delete session %s: %v", sessionID, err))
		}
	}

	// Xóa user sessions set
	err = utils.RedisClient.Del(ctx, userSessionsKey).Err()
	if err != nil {
		logger.Warn("repository", "DeleteUserSessions", fmt.Sprintf("Failed to delete user sessions set: %v", err))
	}

	logger.Info("repository", "DeleteUserSessions", fmt.Sprintf("All sessions deleted for user: %s", userID))
	return nil
}

// IsSessionActive kiểm tra xem session có active không
func (r *sessionRedisRepository) IsSessionActive(ctx context.Context, accessToken string) (bool, error) {
	session, err := r.GetSessionByToken(ctx, accessToken)
	if err != nil {
		return false, nil // Session không tồn tại trong Redis
	}

	// Kiểm tra session có active và chưa hết hạn không
	if !session.IsActive {
		return false, nil
	}

	if time.Now().After(session.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

