package repositories

import (
	"chat-app/internal/core/domain"
	"time"
)

type SessionRepository interface {
	CreateSession(session domain.Session, ttl time.Duration) error
	AddChatToSession(sessionID string, chatID string, role string) error
	RemoveChatFromSession(sessionID string, chatID string) error
	UpdateChatRole(sessionID string, chatID string, role string) error
	IsUserInChat(sessionID string, chatID string) (string, error)
	DeleteSession(sessionID string) error
}
