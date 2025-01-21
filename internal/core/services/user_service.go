package services

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"errors"
	"fmt"
)

type UserService struct {
	User repositories.UserRepository
}

func NewUserService(user repositories.UserRepository) *UserService {
	return &UserService{User: user}
}

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

func (us *UserService) Login(username, password string) (domain.ID, error) {
	if username == "" {
		return "", fmt.Errorf("username is required")
	}
	if password == "" {
		return "", fmt.Errorf("password is required")
	}

	userID, err := us.User.Login(username, password)
	if err != nil {
		if errors.Is(err, repositories.ErrWrongLoginInfo) {
			return "", fmt.Errorf("wrong login info: %v", err)
		}
		return "", fmt.Errorf("failed to login:%v", err)
	}
	return userID, nil
}

func (us *UserService) GetChatIDList(userID domain.ID) ([]domain.ID, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	chatIDList, err := us.User.GetChatIDList(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat list:%v", err)
	}
	return chatIDList, nil
}

func (us *UserService) GetUserInfo(userID domain.ID) (domain.User, error) {
	user, err := us.User.GetUserInfo(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return domain.User{}, fmt.Errorf("%w:%v", repositories.ErrUserNotFound, err)
		}
		return domain.User{}, fmt.Errorf("failed to get user info:%v", err)
	}
	return user, nil
}

func (us *UserService) AddContact(userID, contactID domain.ID) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if contactID == "" {
		return fmt.Errorf("contact ID is required")
	}

	err := us.User.AddContact(userID, contactID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("failed to add contact: %w", err)
	}
	return nil
}

func (us *UserService) RemoveContact(userID, contactID domain.ID) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if contactID == "" {
		return fmt.Errorf("contact ID is required")
	}

	err := us.User.RemoveContact(userID, contactID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("failed to remove contact: %w", err)
	}
	return nil
}
