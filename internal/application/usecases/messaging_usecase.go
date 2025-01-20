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
	isUserInChat, err := m.SessionService.IsUserInChat(sessionID, chatID)

	if err != nil {
		if errors.Is(err, repositories.ErrChatNotFound) {
			return fmt.Errorf("user is not in chat or chat ID is wrong: %v", err)
		}
		return err
	}
	if isUserInChat == "" {
		return fmt.Errorf("user role in chat is missing")
	}

	session, err := m.SessionService.GetSession(sessionID)
	if err != nil {
		return err
	}

	if err = m.MessageService.SendMessage(chatID, session.UserID, message); err != nil {
		return err
	}
	return nil

}

func (m *Messaging) DeleteMessage(chatID domain.ID, sessionID domain.ID) error {
	isUserInChat, err := m.SessionService.IsUserInChat(sessionID, chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrChatNotFound) {
			return fmt.Errorf("user is not in chat or chat ID is wrong: %v", err)
		}
		return err
	}

	if isUserInChat == "" {
		return fmt.Errorf("user role in chat is missing")
	}

	if isUserInChat == "owner" || isUserInChat == "admin" || isUserInChat == "user" {
		if err = m.SessionService.DeleteSession(sessionID); err != nil {
			return err
		}
	}
	return nil
}
