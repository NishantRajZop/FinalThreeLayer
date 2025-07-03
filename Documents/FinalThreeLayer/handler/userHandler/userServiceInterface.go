package userHandler

import (
	"FinalThreeLayer/models"
)

type userService interface {
	CreateUser(models.User) error
	GetUserByID(id int) (models.User, error)
}
