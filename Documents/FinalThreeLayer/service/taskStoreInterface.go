package service

import (
	"FinalThreeLayer/models"
)

type TaskStore interface {
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id int) (models.Task, error)
	CreateTask(task models.Task) error
	UpdateTask(id int) error
	DeleteTask(id int) error
}
