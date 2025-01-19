package services

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"fmt"
)

type UserService struct {
	User repositories.UserRepository
}

//type UserRepository interface {
//	Register(user domain.User) (userID domain.ID, err error)
//	Login(username, password string) (isAuthorized bool, err error)
//	GetChatList(userID domain.ID) (chatList []string, err error)
//}

func NewUserService(user repositories.UserRepository) *UserService {
	return &UserService{User: user}
}

//type User struct {
//	ID          ID
//	Username    string
//	FirstName   string
//	LastName    string
//	Password    string
//	Gender      Gender
//	Email       string
//	Contacts    []ID
//	DateOfBirth time.Time
//	CreatedTime time.Time
//	DeletedTime time.Time
//}

func ValidateUser(user domain.User) error {
	if user.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}

	if user.Gender == 0 {
		return fmt.Errorf("gender is required")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	if user.DateOfBirth == nil {
		return fmt.Errorf("date of birth is required")
	}
	return nil

}

func (us *UserService) Register(user domain.User) (userID domain.ID, err error) {
	if err = ValidateUser(user); err != nil {
		return "", fmt.Errorf("missing user fields:%v", err)
	}

	userID, err = us.User.Register(user)
	if err != nil {
		return "", fmt.Errorf("failed to register user:%v", err)
	}

	return userID, nil
}

func (us *UserService) Login(username, password string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("username is required")
	}
	if password == "" {
		return false, fmt.Errorf("password is required")
	}

	isAuthenticated, err := us.User.Login(username, password)
	if err != nil {
		return false, fmt.Errorf("failed to login:%v", err)
	}
	return isAuthenticated, nil
}

func (us *UserService) GetChatList(userID domain.ID) (chatList []string, err error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	chatList, err = us.User.GetChatList(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat list:%v", err)
	}
	return chatList, nil
}
