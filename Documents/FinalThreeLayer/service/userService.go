package service

import (
	"FinalThreeLayer/models"
	"errors"
)

type userService struct {
	store UserStore
}

func NewUserService(store UserStore) *userService {
	return &userService{store: store}
}

func (uService *userService) CreateUser(user models.User) error {
	// you can implement all your Validations here
	if user.Name == "" {
		return errors.New("user name cannot be empty")
	}
	if user.ID == 0 {
		return errors.New("invalid user ID")
	}
	return uService.store.CreateUser(user)
}

func (uService *userService) GetUserByID(id int) (models.User, error) {
	if id == 0 {
		return models.User{}, errors.New("invalid user ID")
	}
	return uService.store.GetUserByID(id)
}
