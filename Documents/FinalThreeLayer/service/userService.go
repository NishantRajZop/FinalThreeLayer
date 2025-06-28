package service

import (
	"FinalThreeLayer/models"
)

type UserStore interface {
	CreateUser(models.User) error
	GetUserByID(id int) (models.User, error)
}

type userService struct {
	store UserStore
}

func NewUserService(store UserStore) *userService {
	return &userService{store: store}
}

func (uService *userService) CreateUser(user models.User) error {
	// you can implement all your Validations here
	return uService.store.CreateUser(user)
}

func (uService *userService) GetUserByID(id int) (models.User, error) {
	// you can implement all Your validations here
	return uService.store.GetUserByID(id)
}
