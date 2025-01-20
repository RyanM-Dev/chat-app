package usecases

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/services"
	"fmt"
)

type ChatManagement struct {
	ChatService    *services.ChatService
	SessionService *services.SessionService
}

func NewChatManagement(chatService *services.ChatService, sessionService *services.SessionService) *ChatManagement {
	return &ChatManagement{
		ChatService:    chatService,
		SessionService: sessionService,
	}
}

func (cm *ChatManagement) CreateChat(chat domain.Chat, sessionID domain.ID) error {
	session, err := cm.SessionService.GetSession(sessionID)
	if err != nil {
		return err
	}

	_, exist := session.ChatNameAndID[chat.Name]
	if exist {
		return fmt.Errorf("chat name %s already exists", chat.Name)
	}

	chatID, err := cm.ChatService.CreateChat(chat)

	if err != nil {
		return err
	}

	err = cm.SessionService.AddChatToSession(sessionID, chatID, "owner")
	if err != nil {
		return err
	}

	return nil
}

func (cm *ChatManagement) FindChat(chatID, sessionID domain.ID) (domain.Chat, error) {
	_, err := cm.chatAuthorization(domain.ID(chatID), sessionID)
	if err != nil {
		return domain.Chat{}, err
	}

	chat, err := cm.ChatService.FindChat(chatID)
	if err != nil {
		return domain.Chat{}, err
	}
	return chat, nil
}

func (cm *ChatManagement) UpdateChatName(currentChatName, nextChatName string, sessionID domain.ID) error {
	session, err := cm.SessionService.GetSession(sessionID)
	if err != nil {
		return err
	}

	chatID, exist := session.ChatNameAndID[currentChatName]
	if !exist {
		return fmt.Errorf("wrong chat name")
	}

	chat, err := cm.ChatService.FindChat(domain.ID(chatID))
	if err != nil {
		return err
	}

	userRole, err := cm.chatAuthorization(domain.ID(chatID), sessionID)
	if err != nil {
		return err
	}

	if userRole != domain.Owner {
		return fmt.Errorf("only the owner of this chat can change the name of %s", currentChatName)
	}

	chat.Name = nextChatName

	err = cm.ChatService.UpdateChatName(chat)
	return nil
}

func (cm *ChatManagement) DeleteChat(chatID, sessionID domain.ID) error {
	userRole, err := cm.chatAuthorization(chatID, sessionID)
	if err != nil {
		return err
	}
	if userRole != domain.Owner {
		return fmt.Errorf("you are not authorized to delete this chat")
	}
	err = cm.ChatService.DeleteChat(chatID)
	if err != nil {
		return err
	}

	err = cm.SessionService.RemoveChatFromSession(sessionID, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (cm *ChatManagement) GetMessages(chatName string, sessionID domain.ID) ([]domain.Message, error) {
	session, err := cm.SessionService.GetSession(sessionID)
	if err != nil {
		return nil, err

	}
	chatID, exist := session.ChatNameAndID[chatName]
	if !exist {
		return nil, fmt.Errorf("user is not in this chat or chat doesn't exist: %v", chatName)
	}

	messages, err := cm.ChatService.GetMessages(domain.ID(chatID))
	if err != nil {
		return nil, err
	}
	return messages, nil

}

func (cm *ChatManagement) AddUser(chatID, sessionID domain.ID, userIDs []domain.ID) error {
	userRole, err := cm.chatAuthorization(chatID, sessionID)
	if err != nil {
		return err
	}
	if userRole != domain.Admin && userRole != domain.Owner {
		return fmt.Errorf("you are not authorized to remove from this chat")
	}

	err = cm.ChatService.AddUser(chatID, userIDs)
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		session, err := cm.SessionService.GetSessionByUserID(userID)
		if err != nil {
			return err
		}
		err = cm.SessionService.AddChatToSession(session.SessionID, chatID, domain.Normal)
		if err != nil {
			return err
		}

	}

	return nil

}

func (cm *ChatManagement) RemoveUser(chatID, sessionID domain.ID, userIDs []domain.ID) error {
	userRole, err := cm.chatAuthorization(chatID, sessionID)
	if err != nil {
		return err
	}

	if userRole != domain.Admin && userRole != domain.Owner {
		return fmt.Errorf("you are not authorized to remove from this chat")
	}

	err = cm.ChatService.RemoveUser(chatID, userIDs)
	if err != nil {
		return err
	}
	return nil

}

func (cm *ChatManagement) GetMembers(chatID, sessionID domain.ID) ([]domain.ID, error) {
	_, err := cm.chatAuthorization(chatID, sessionID)
	if err != nil {
		return nil, err
	}

	chatMembers, err := cm.ChatService.GetMembers(chatID)
	if err != nil {
		return nil, err
	}
	return chatMembers, nil

}

func (cm *ChatManagement) SetAdmin(chatID, sessionID domain.ID, userIDs []domain.ID) error {
	role, err := cm.chatAuthorization(chatID, sessionID)
	if err != nil {
		return err
	}

	if role != domain.Admin && role != domain.Owner {
		return fmt.Errorf("you are not authorized to set admin in this chat")
	}

	for _, userID := range userIDs {
		err = cm.ChatService.SetAdmin(userID, chatID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cm *ChatManagement) chatAuthorization(chatID, sessionID domain.ID) (string, error) {
	session, err := cm.SessionService.GetSession(sessionID)
	if err != nil {
		return "", err

	}
	chat, err := cm.ChatService.FindChat(chatID)
	if err != nil {
		return "", err
	}

	userChatID, exist := session.ChatNameAndID[chat.Name]
	if !exist {
		return "", fmt.Errorf("user is not in this chat or chat doesn't exist: %v", chat.Name)
	}

	if domain.ID(userChatID) != chatID {
		return "", fmt.Errorf("chatID is not valid: %v", chatID)
	}

	for _, admin := range chat.Admins {
		if admin == session.UserID {
			role := domain.Admin
			return role, nil
		}
	}

	if chat.Owner == session.UserID {
		role := domain.Owner
		return role, nil
	}

	role := domain.Normal
	return role, nil

}
