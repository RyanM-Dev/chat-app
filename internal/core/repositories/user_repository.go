package repositories

import (
	"chat-app/internal/core/domain"
	"errors"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrWrongLoginInfo = errors.New("wrong login info")
)

type UserRepository interface {
	Register(user domain.User) (userID domain.ID, err error)
	Login(username, password string) (userID domain.ID, err error)
	GetChatIDList(userID domain.ID) (chatIDList []string, err error)
	GetUserInfo(userID domain.ID) (user domain.User, err error)
}
