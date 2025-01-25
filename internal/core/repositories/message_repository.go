package repositories

import "chat-app/internal/core/domain"

type MessageRepository interface {
	SendMessage(chatID, userID domain.ID, message string) error
	DeleteMessage(messageID domain.ID) error
	GetMessage(messageID domain.ID) (domain.Message, error)
}
