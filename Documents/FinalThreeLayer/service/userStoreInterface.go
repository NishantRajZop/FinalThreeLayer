package service

import (
	"FinalThreeLayer/models"
)

type UserStore interface {
	CreateUser(models.User) error
	GetUserByID(id int) (models.User, error)
}
