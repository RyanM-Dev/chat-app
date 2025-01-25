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

func (ms *MessageService) GetMessage(messageID domain.ID) (domain.Message, error) {
	if messageID == "" {
		return domain.Message{}, fmt.Errorf("messageID cannot be empty")
	}
	message, err := ms.Message.GetMessage(messageID)
	if err != nil {
		return domain.Message{}, fmt.Errorf("error getting message: %w", err)
	}

	return message, nil
}

func NewMessageService(message repositories.MessageRepository) *MessageService {
	return &MessageService{
		Message: message,
	}
}

func (ms *MessageService) SendMessage(chatID, userID domain.ID, message string) error {
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	if err := ms.Message.SendMessage(chatID, userID, message); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func (ms *MessageService) DeleteMessage(messageID domain.ID) error {
	if messageID == "" {
		return fmt.Errorf("messageID cannot be empty")
	}

	if err := ms.Message.DeleteMessage(messageID); err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}
