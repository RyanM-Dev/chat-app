package repositories

import "chat-app/internal/core/domain"

type UserRepository interface {
	Register(user domain.User) (userID domain.ID, err error)
	Login(username, password string) (isAuthorized bool, err error)
	GetChatList(userID domain.ID) (chatList []string, err error)
}
