package services

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"errors"
	"fmt"
	"time"
)

type SessionService struct {
	SessionRepo repositories.SessionRepository
}

func NewSessionService(sessionRepo repositories.SessionRepository) *SessionService {
	return &SessionService{SessionRepo: sessionRepo}
}

func (s *SessionService) CreateSession(session domain.Session, ttl time.Duration) error {
	if session.SessionID == "" {
		return fmt.Errorf("sessionID cannot be empty")
	}
	if session.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	return s.SessionRepo.CreateSession(session, ttl)
}

func (s *SessionService) AddChat(sessionID string, chatID string, role string) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID cannot be empty")
	}
	if chatID == "" {
		return fmt.Errorf("chatID cannot be empty")
	}
	if role == "" {
		return fmt.Errorf("role cannot be empty")
	}
	return s.SessionRepo.AddChatToSession(sessionID, chatID, role)
}

func (s *SessionService) RemoveChat(sessionID string, chatID string) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID cannot be empty")
	}
	if chatID == "" {
		return fmt.Errorf("chatID cannot be empty")
	}
	if err := s.SessionRepo.RemoveChatFromSession(sessionID, chatID); err != nil {
		return err
	}
	return nil
}

func (s *SessionService) UpdateChatRole(sessionID string, chatID string, role string) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID cannot be empty")
	}
	if chatID == "" {
		return fmt.Errorf("chatID cannot be empty")
	}
	if role == "" {
		return fmt.Errorf("role cannot be empty")
	}
	if err := s.SessionRepo.UpdateChatRole(sessionID, chatID, role); err != nil {
		return err
	}
	return nil
}

func (s *SessionService) IsUserInChat(sessionID string, chatID string) (string, error) {
	if sessionID == "" {
		return "", fmt.Errorf("sessionID cannot be empty")
	}
	if chatID == "" {
		return "", fmt.Errorf("chatID cannot be empty")
	}
	role, err := s.SessionRepo.IsUserInChat(sessionID, chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return "", fmt.Errorf("user not found in chat:%v", err)
		}
		return "", fmt.Errorf("error while checking if user is in chat:%v", err)
	}
	return role, nil
}

func (s *SessionService) DeleteSession(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID cannot be empty")
	}
	if err := s.SessionRepo.DeleteSession(sessionID); err != nil {
		return err
	}
	return nil
}
