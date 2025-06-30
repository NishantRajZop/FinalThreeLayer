package service

import (
	"FinalThreeLayer/models"
	"errors"
	"testing"
)

// MockTaskStore implements TaskStore for testing
type MockTaskStore struct {
	GetAllTasksFunc func() ([]models.Task, error)
	GetTaskByIDFunc func(id int) (models.Task, error)
	CreateTaskFunc  func(task models.Task) error
	UpdateTaskFunc  func(id int) error
	DeleteTaskFunc  func(id int) error
}

func (m *MockTaskStore) GetAllTasks() ([]models.Task, error) {
	return m.GetAllTasksFunc()
}

func (m *MockTaskStore) GetTaskByID(id int) (models.Task, error) {
	return m.GetTaskByIDFunc(id)
}

func (m *MockTaskStore) CreateTask(task models.Task) error {
	return m.CreateTaskFunc(task)
}

func (m *MockTaskStore) UpdateTask(id int) error {
	return m.UpdateTaskFunc(id)
}

func (m *MockTaskStore) DeleteTask(id int) error {
	return m.DeleteTaskFunc(id)
}

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name        string
		mockStore   *MockTaskStore
		wantTasks   []models.Task
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - get all tasks",
			mockStore: &MockTaskStore{
				GetAllTasksFunc: func() ([]models.Task, error) {
					return []models.Task{
						{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
						{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
					}, nil
				},
			},
			wantTasks: []models.Task{
				{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
				{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
			},
			wantErr: false,
		},
		{
			name: "error - database failure",
			mockStore: &MockTaskStore{
				GetAllTasksFunc: func() ([]models.Task, error) {
					return nil, errors.New("database error")
				},
			},
			wantTasks:   nil,
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTaskService(tt.mockStore)
			tasks, err := service.GetAllTasks()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.expectedErr.Error() {
				t.Errorf("GetAllTasks() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if len(tasks) != len(tt.wantTasks) {
				t.Errorf("GetAllTasks() returned %d tasks, want %d", len(tasks), len(tt.wantTasks))
				return
			}

			for i := range tasks {
				if tasks[i] != tt.wantTasks[i] {
					t.Errorf("GetAllTasks() task %d = %v, want %v", i, tasks[i], tt.wantTasks[i])
				}
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name        string
		taskID      int
		mockStore   *MockTaskStore
		wantTask    models.Task
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - get task by ID",
			taskID: 1,
			mockStore: &MockTaskStore{
				GetTaskByIDFunc: func(id int) (models.Task, error) {
					return models.Task{ID: id, Title: "Task 1", UserID: 1}, nil
				},
			},
			wantTask: models.Task{ID: 1, Title: "Task 1", UserID: 1},
			wantErr:  false,
		},
		{
			name:   "error - task not found",
			taskID: 999,
			mockStore: &MockTaskStore{
				GetTaskByIDFunc: func(id int) (models.Task, error) {
					return models.Task{}, errors.New("task not found")
				},
			},
			wantTask:    models.Task{},
			wantErr:     true,
			expectedErr: errors.New("task not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTaskService(tt.mockStore)
			task, err := service.GetTaskByID(tt.taskID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.expectedErr.Error() {
				t.Errorf("GetTaskByID() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if task != tt.wantTask {
				t.Errorf("GetTaskByID() = %v, want %v", task, tt.wantTask)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name        string
		task        models.Task
		mockStore   *MockTaskStore
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - create task",
			task: models.Task{Title: "New Task", UserID: 1},
			mockStore: &MockTaskStore{
				CreateTaskFunc: func(task models.Task) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "error - database failure",
			task: models.Task{Title: "New Task", UserID: 1},
			mockStore: &MockTaskStore{
				CreateTaskFunc: func(task models.Task) error {
					return errors.New("database error")
				},
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name: "error - empty title",
			task: models.Task{Title: "", UserID: 1},
			mockStore: &MockTaskStore{
				CreateTaskFunc: func(task models.Task) error {
					return errors.New("client Side error")
				},
			},
			wantErr:     true,
			expectedErr: errors.New("task title cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTaskService(tt.mockStore)
			err := service.CreateTask(tt.task)

			if err != nil && tt.wantErr == false {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.wantErr == true {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		taskID      int
		mockStore   *MockTaskStore
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - update task",
			taskID: 1,
			mockStore: &MockTaskStore{
				UpdateTaskFunc: func(id int) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name:   "error - database failure",
			taskID: 1,
			mockStore: &MockTaskStore{
				UpdateTaskFunc: func(id int) error {
					return errors.New("database error")
				},
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		//{
		//	name:   "error - invalid task ID",
		//	taskID: 0,
		//	mockStore: &MockTaskStore{
		//		UpdateTaskFunc: func(id int) error {
		//			return nil // Shouldn't be called
		//			// it should not be called because , this is validation and it should be done service layer
		//			// validation section, not here .It prevents unnecessary database calls
		//		},
		//	},
		//	wantErr:     true,
		//	expectedErr: errors.New("invalid task ID"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTaskService(tt.mockStore)
			err := service.UpdateTask(tt.taskID)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("UpdateTask() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name        string
		taskID      int
		mockStore   *MockTaskStore
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - delete task",
			taskID: 1,
			mockStore: &MockTaskStore{
				DeleteTaskFunc: func(id int) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name:   "error - database failure",
			taskID: 1,
			mockStore: &MockTaskStore{
				DeleteTaskFunc: func(id int) error {
					return errors.New("database error")
				},
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		//{
		//	name:   "error - invalid task ID",
		//	taskID: 0,
		//	mockStore: &MockTaskStore{
		//		DeleteTaskFunc: func(id int) error {
		//			return nil // Shouldn't be called
		//		},
		//	},
		//	wantErr:     true,
		//	expectedErr: errors.New("invalid task ID"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTaskService(tt.mockStore)
			err := service.DeleteTask(tt.taskID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("DeleteTask() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
