package repositories

import "chat-app/internal/core/domain"

type MessageRepository interface {
	SendMessage(chatID domain.ID, message string) error
	DeleteMessage(chatID, messageID domain.ID) error
}
