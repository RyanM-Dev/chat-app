package usecases

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"chat-app/internal/core/services"
	"errors"
	"fmt"
)

type Messaging struct {
	ChatService    *services.ChatService
	MessageService *services.MessageService
	SessionService *services.SessionService
}

func NewMessaging(chatService *services.ChatService, messageService *services.MessageService, sessionService *services.SessionService) *Messaging {
	return &Messaging{
		ChatService:    chatService,
		MessageService: messageService,
		SessionService: sessionService,
	}
}

func (m *Messaging) SendMessage(chatID domain.ID, sessionID domain.ID, message string) error {

	chat, err := m.ChatService.FindChat(chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrChatNotFound) {
			return fmt.Errorf("chat doesn't exist :%v", err)
		}
		return err
	}

}
