package repositories

import (
	"chat-app/internal/core/domain"
	"errors"
)

var (
	ErrChatNotFound         = errors.New("chat not found")
	ErrMissingChatParameter = errors.New("missing chat parameters")
	ErrDuplicateChat        = errors.New("duplicate chat")
)

type ChatRepository interface {
	CreateChat(chat domain.Chat) (chatID domain.ID, err error)
	FindChat(chatID domain.ID) (chat domain.Chat, err error)
	UpdateChat(chat domain.Chat) error
	DeleteChat(chatID domain.ID) error
	GetMessages(chatID domain.ID) ([]domain.Message, error)
	AddUser(chatID domain.ID, userID []domain.ID) error
	RemoveUser(chatID domain.ID, userID []domain.ID) error
	GetMembers(chatID domain.ID) ([]domain.ID, error)
	SetAdmin(userID, chatID domain.ID) error
}
