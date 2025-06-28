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

type taskService struct {
	store TaskStore
}

func NewTaskService(store TaskStore) *taskService {
	return &taskService{store: store}
}

func (tService *taskService) GetAllTasks() ([]models.Task, error) {
	//You can Perform All Your Validations here
	return tService.store.GetAllTasks()
}

func (tService *taskService) GetTaskByID(id int) (models.Task, error) {
	//You can Perform All Your Validations here
	return tService.store.GetTaskByID(id)
}

func (tService *taskService) CreateTask(task models.Task) error {
	//You can Perform All Your Validations here
	return tService.store.CreateTask(task)
}

func (tService *taskService) UpdateTask(id int) error {
	//You can Perform All Your Validations here
	return tService.store.UpdateTask(id)
}

func (tService *taskService) DeleteTask(id int) error {
	//You can Perform All Your Validations here
	return tService.store.DeleteTask(id)
}
