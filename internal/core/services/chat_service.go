package services

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"errors"
	"fmt"
)

type ChatService struct {
	Chat repositories.ChatRepository
}

//type ChatRepository interface {
//	CreateChat(chat domain.Chat) (chatID domain.ID, err error)
//	FindChat(chatID domain.ID) (chat domain.Chat, err error)
//	UpdateChatName(chat domain.Chat) error
//	DeleteChat(chatID domain.ID) error
//	GetMessages(chatID domain.ID) ([]domain.Message, error)
//	AddUser(chatID domain.ID, userID domain.ID) error
//	RemoveUser(chatID domain.ID, userID domain.ID) error
//	GetMembers(chatID domain.ID) ([]domain.ID, error)
//}

func ValidateChat(chat domain.Chat) error {
	if chat.ID == "" {
		return fmt.Errorf("%w: Chat.ID", repositories.ErrMissingChatParameter)
	}
	if chat.Name == "" {
		return fmt.Errorf("%w: Chat.Name", repositories.ErrMissingChatParameter)
	}
	if chat.Owner == "" {
		return fmt.Errorf("%w: Chat.Owner", repositories.ErrMissingChatParameter)
	}
	if chat.Members == nil {
		return fmt.Errorf("%w: Chat.Members", repositories.ErrMissingChatParameter)
	}
	if chat.CreatedTime == nil {
		return fmt.Errorf("%w: Chat.CreatedTime", repositories.ErrMissingChatParameter)
	}
	if chat.ChatType == 0 {
		return fmt.Errorf("%w: Chat.ChatType", repositories.ErrMissingChatParameter)
	}
	return nil
}

func NewChatService(chat repositories.ChatRepository) *ChatService {
	return &ChatService{chat}
}

func (cs *ChatService) FindChat(chatID domain.ID) (domain.Chat, error) {
	chat, err := cs.Chat.FindChat(chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrChatNotFound) {
			return domain.Chat{}, fmt.Errorf(" chat doesn't exist: %w", repositories.ErrChatNotFound)
		}
		return domain.Chat{}, fmt.Errorf(" failed to find the chat: %v", err)
	}

	err = ValidateChat(chat)
	if err != nil {
		return domain.Chat{}, err
	}

	return chat, nil
}

func (cs *ChatService) CreateChat(chat domain.Chat) (domain.ID, error) {
	if err := ValidateChat(chat); err != nil {
		return "", err
	}

	chatID, err := cs.Chat.CreateChat(chat)
	if err != nil {

		if errors.Is(err, repositories.ErrDuplicateChat) {
			return "", fmt.Errorf("chat already exists: %v", err)
		}

		return "", fmt.Errorf("falied to create Chat: %v", err)
	}

	return chatID, nil

}

func (cs *ChatService) UpdateChatName(UpdateChat domain.Chat) error {
	if err := ValidateChat(UpdateChat); err != nil {
		return err
	}

	chat, err := cs.Chat.FindChat(UpdateChat.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrChatNotFound) {
			return fmt.Errorf("chat for update not found: %w", repositories.ErrChatNotFound)
		}
		return fmt.Errorf("falied to update Chat: %v", err)
	}

	chat.Name = UpdateChat.Name
	chat.Owner = UpdateChat.Owner

	err = cs.Chat.UpdateChatName(chat)
	if err != nil {
		return fmt.Errorf("falied to update Chat: %v", err)
	}
	return nil
}

func (cs *ChatService) DeleteChat(chatID domain.ID) error {
	if chatID == "" {
		return fmt.Errorf("missing chat id")
	}
	//_, err := cs.Chat.FindChat(chatID)
	//if err != nil {
	//	if errors.Is(err, repositories.ErrChatNotFound) {
	//		return fmt.Errorf("chat for delete not found: %w", repositories.ErrChatNotFound)
	//	}
	//	return fmt.Errorf("failed to delete Chat: %v", err)
	//}
	err := cs.Chat.DeleteChat(chatID)
	if err != nil {
		return fmt.Errorf("falied to delete Chat: %v", err)
	}
	return nil
}

func (cs *ChatService) GetMessages(chatID domain.ID) ([]domain.Message, error) {
	if chatID == "" {
		return nil, fmt.Errorf("missing chat id")
	}
	//chat, err := cs.Chat.FindChat(chatID)
	//if err != nil {
	//	if errors.Is(err, repositories.ErrChatNotFound) {
	//		return nil, fmt.Errorf("chat for get messages not found: %w", repositories.ErrChatNotFound)
	//	}
	//	return nil, fmt.Errorf("failed to find chat to get messages: %v", err)
	//}

	messages, err := cs.Chat.GetMessages(chatID)
	if err != nil {
		return nil, fmt.Errorf("falied to get messages: %v", err)
	}
	return messages, nil
}

func (cs *ChatService) AddUser(chatID domain.ID, userIDs []domain.ID) error {
	if chatID == "" {
		return fmt.Errorf("missing chat id")
	}
	if userIDs == nil {
		return fmt.Errorf("missing user ids")
	}
	//chat, err := cs.Chat.FindChat(chatID)
	//if err != nil {
	//	if errors.Is(err, repositories.ErrChatNotFound) {
	//		return fmt.Errorf("chat for adding user not found: %w", repositories.ErrChatNotFound)
	//	}
	//	return fmt.Errorf("failed to find chat to add user: %v", err)
	//}
	//
	//for _, userID := range userIDs {
	//	for _, member := range chat.Members {
	//		if member == userID {
	//			return fmt.Errorf("user %s is already member of chat", userID)
	//		}
	//	}
	//}

	err := cs.Chat.AddUser(chatID, userIDs)
	if err != nil {
		return fmt.Errorf("falied to add user to chat: %v", err)
	}

	return nil
}

func (cs *ChatService) RemoveUser(chatID domain.ID, userIDs []domain.ID) error {
	if chatID == "" {
		return fmt.Errorf("missing chat id")
	}
	if userIDs == nil {
		return fmt.Errorf("missing user ids")
	}
	//_, err := cs.Chat.FindChat(chatID)
	//if err != nil {
	//	if errors.Is(err, repositories.ErrChatNotFound) {
	//		return fmt.Errorf("chat for removing user not found: %w", repositories.ErrChatNotFound)
	//	}
	//	return fmt.Errorf("failed to find chat to remove user: %v", err)
	//}

	err := cs.Chat.RemoveUser(chatID, userIDs)
	if err != nil {
		return fmt.Errorf("falied to remove users from chat: %v", err)
	}
	return nil
}

func (cs *ChatService) GetMembers(chatID domain.ID) ([]domain.ID, error) {
	if chatID == "" {
		return nil, fmt.Errorf("missing chat id")
	}

	members, err := cs.Chat.GetMembers(chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrChatNotFound) {
			return nil, fmt.Errorf("%w: %v", repositories.ErrChatNotFound, err)
		}

		return nil, fmt.Errorf("falied to get members: %v", err)
	}

	return members, nil
}
func (cs *ChatService) SetAdmin(adminID, chatID domain.ID) error {
	if adminID == "" {
		return fmt.Errorf("user ID is empty")
	}
	if chatID == "" {
		return fmt.Errorf("chat ID is empty")
	}
	err := cs.Chat.SetAdmin(adminID, chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user  not found")
		}
		if errors.Is(err, repositories.ErrChatNotFound) {
			return fmt.Errorf("chat  not found")
		}
		return fmt.Errorf("failed to set admin: %v", err)
	}
	return nil
}
