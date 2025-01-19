package services

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"fmt"
)

type MessageService struct {
	Message repositories.MessageRepository
	Chat    repositories.ChatRepository
}

func NewMessageService(message repositories.MessageRepository) *MessageService {
	return &MessageService{
		Message: message,
	}
}

func (ms *MessageService) SendMessage(chatID domain.ID, message string) error {
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	if err := ms.Message.SendMessage(chatID, message); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func (ms *MessageService) DeleteMessage(chatID, messageID domain.ID) error {
	if chatID == "" {
		return fmt.Errorf("chatID cannot be empty")
	}
	if messageID == "" {
		return fmt.Errorf("messageID cannot be empty")
	}
	if err := ms.Message.DeleteMessage(chatID, messageID); err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}
