package usecases

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/services"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type UserManagement struct {
	UserService    *services.UserService
	ChatService    *services.ChatService
	SessionService *services.SessionService
}

func NewUserManagement(userService *services.UserService, sessionService *services.SessionService) *UserManagement {
	return &UserManagement{
		UserService:    userService,
		SessionService: sessionService,
	}
}

// Register creates a new user account by delegating to the UserService, initializes a session, and appends it to Sessions.
func (um *UserManagement) Register(user domain.User) (domain.Session, error) {
	userID, err := um.UserService.Register(user)
	if err != nil {
		return domain.Session{}, err
	}

	newSession := um.createSession(userID)

	return newSession, nil
}

// createSession is a helper function to create a new session for a user.
func (um *UserManagement) createSession(userID domain.ID) domain.Session {
	newSessionID := domain.ID(uuid.New().String())
	return domain.Session{
		SessionID: newSessionID,
		UserID:    userID,
	}
}
func (um *UserManagement) Login(username string, password string) (domain.Session, error) {
	userID, err := um.UserService.Login(username, password)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to log in: %w", err)
	}

	chatIDList, err := um.UserService.GetChatIDList(userID)
	if err != nil {
		return domain.Session{}, err
	}
	if len(chatIDList) == 0 {
		return domain.Session{}, fmt.Errorf("no chats found for user")
	}

	var chatNames []string
	chatMapping := make(domain.ChatMapping)
	for _, chatID := range chatIDList {
		chat, err := um.ChatService.FindChat(chatID)
		if err != nil {
			// Log the error and continue
			log.Printf("failed to find chat for ID %s: %v", chatID, err)
			continue
		}
		chatMapping[chat.Name] = chat.ID
		chatNames = append(chatNames, chat.Name)
	}

	sessionIDUUID := uuid.New()
	if sessionIDUUID == uuid.Nil {
		return domain.Session{}, fmt.Errorf("failed to generate session ID")
	}
	sessionID := domain.ID(sessionIDUUID.String())
	session := domain.Session{
		SessionID:    sessionID,
		UserID:       userID,
		ChatMappings: chatMapping,
		ChatNames:    chatNames,
	}

	return session, nil
}

func (um *UserManagement) AddContact(userID, contactID domain.ID) error {
	// Update the user in the database
	err := um.UserService.AddContact(userID, contactID)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserManagement) RemoveContact(userID, contactID domain.ID) error {
	// Update the user in the database
	err := um.UserService.RemoveContact(userID, contactID)
	if err != nil {
		return err
	}
	return nil
}
