package usecases

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/services"
	"github.com/google/uuid"
)

type UserManagement struct {
	UserService    *services.UserService
	ChatService    *services.ChatService
	SessionService *services.SessionService
	Sessions       *[]domain.Session
}

//type UserRepository interface {
//	Register(user domain.User) (userID domain.ID, err error)
//	Login(username, password string) (isAuthorized bool, err error)
//	GetChatIDList(userID domain.ID) (chatIDList []string, err error)
//	GetUserInfo(userID domain.ID) (user domain.User, err error)
//}

func NewUserManagement(userService *services.UserService, sessionService *services.SessionService, sessions *[]domain.Session) *UserManagement {
	return &UserManagement{
		UserService:    userService,
		SessionService: sessionService,
		Sessions:       sessions,
	}
}

func (um *UserManagement) Register(user domain.User) (domain.Session, error) {
	userID, err := um.UserService.Register(user)
	if err != nil {
		return domain.Session{}, err
	}
	sessionID := domain.ID(uuid.New().String())
	session := domain.Session{
		SessionID: sessionID,
		UserID:    userID,
	}
	*um.Sessions = append(*um.Sessions, session)

	return session, nil

}
func (um *UserManagement) Login(username string, password string) (domain.Session, error) {
	userID, err := um.UserService.Login(username, password)
	if err != nil {
		return domain.Session{}, err
	}

	chatIDList, err := um.UserService.GetChatIDList(userID)
	if err != nil {
		return domain.Session{}, err
	}

	var chatNameList []string
	var chatNameAndID map[string]string
	for _, chatID := range chatIDList {
		chat, err := um.ChatService.FindChat(domain.ID(chatID))
		if err != nil {
			return domain.Session{}, err
		}
		chatNameAndID[chat.Name] = chatID
		chatNameList = append(chatNameList, chat.Name)
	}
	sessionID := domain.ID(uuid.New().String())
	session := domain.Session{
		SessionID:     sessionID,
		UserID:        userID,
		ChatNameAndID: chatNameAndID,
		ChatNameList:  chatNameList,
	}

	return session, nil
}
