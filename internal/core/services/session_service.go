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

func (s *SessionService) GetSession(sessionId domain.ID) (domain.Session, error) {
	if sessionId == "" {
		return domain.Session{}, errors.New("sessionID cannot be empty")
	}
	session, err := s.SessionRepo.GetSession(sessionId)
	if err != nil {
		return domain.Session{}, fmt.Errorf("error getting session: %v", err)
	}
	return session, nil
}

func (s *SessionService) GetSessionByUserID(userID domain.ID) (domain.Session, error) {
	if userID == "" {
		return domain.Session{}, errors.New("userID cannot be empty")
	}
	session, err := s.SessionRepo.GetSessionByUserID(userID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("error getting session: %v", err)
	}
	return session, nil
}

func (s *SessionService) AddChatToSession(sessionID domain.ID, chatID domain.ID, role string) error {
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

func (s *SessionService) RemoveChatFromSession(sessionID domain.ID, chatID domain.ID) error {
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

func (s *SessionService) UpdateChatRole(sessionID domain.ID, chatID domain.ID, role string) error {
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

func (s *SessionService) IsUserInChat(sessionID domain.ID, chatID domain.ID) (string, error) {
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

func (s *SessionService) DeleteSession(sessionID domain.ID) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID cannot be empty")
	}
	if err := s.SessionRepo.DeleteSession(sessionID); err != nil {
		return err
	}
	return nil
}
