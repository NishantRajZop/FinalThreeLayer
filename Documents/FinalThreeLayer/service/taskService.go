package service

import (
	"FinalThreeLayer/models"
	"errors"
)

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
	if task.Title == "" {
		return errors.New("task title cannot be empty")
	}
	return tService.store.CreateTask(task)
}

func (tService *taskService) UpdateTask(id int) error {
	if id == 0 {
		return errors.New("invalid task ID")
	}
	return tService.store.UpdateTask(id)
}

func (tService *taskService) DeleteTask(id int) error {
	//You can Perform All Your Validations here
	if id == 0 {
		return errors.New("invalid task ID")
	}
	return tService.store.DeleteTask(id)
}
