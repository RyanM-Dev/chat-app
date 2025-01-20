package repositories

import (
	"chat-app/internal/core/domain"
	"time"
)

type SessionRepository interface {
	CreateSession(session domain.Session, ttl time.Duration) error
	GetSession(sessionID domain.ID) (domain.Session, error)
	GetSessionByUserID(userID domain.ID) (domain.Session, error)
	AddChatToSession(sessionID domain.ID, chatID domain.ID, role string) error
	RemoveChatFromSession(sessionID domain.ID, chatID domain.ID) error
	UpdateChatRole(sessionID domain.ID, chatID domain.ID, role string) error
	IsUserInChat(sessionID domain.ID, chatID domain.ID) (role string, err error)
	DeleteSession(sessionID domain.ID) error
}
