package repositories

import (
	"chat-app/internal/core/domain"
)

type ChatRepository interface {
	CreateChat(chat domain.Chat) (chatID domain.ID, err error)
	FindChat(chatID domain.ID) (chat domain.Chat, err error)
	UpdateChatName(chat domain.Chat) error
	DeleteChat(chatID domain.ID) error
	GetMessages(chatID domain.ID) ([]domain.Message, error)
	AddUser(chatID domain.ID, userIDs []domain.ID) error
	RemoveUser(chatID domain.ID, userID []domain.ID) error
	GetMembers(chatID domain.ID) ([]domain.ID, error)
	SetAdmin(adminID, chatID domain.ID) error
}
