package repositories

import "chat-app/internal/core/domain"

type MessageRepository interface {
	SendMessage(chatID, userID domain.ID, message string) error
	DeleteMessage(chatID, userID, messageID domain.ID) error
}
