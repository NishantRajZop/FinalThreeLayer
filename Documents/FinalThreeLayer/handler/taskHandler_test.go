package handler

import (
	"FinalThreeLayer/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockTaskService implements taskService for testing
type MockTaskService struct {
	GetAllTasksFunc func() ([]models.Task, error)
	GetTaskByIDFunc func(id int) (models.Task, error)
	CreateTaskFunc  func(task models.Task) error
	UpdateTaskFunc  func(id int) error
	DeleteTaskFunc  func(id int) error
}

func (m *MockTaskService) GetAllTasks() ([]models.Task, error) {
	return m.GetAllTasksFunc()
}

func (m *MockTaskService) GetTaskByID(id int) (models.Task, error) {
	return m.GetTaskByIDFunc(id)
}

func (m *MockTaskService) CreateTask(task models.Task) error {
	return m.CreateTaskFunc(task)
}

func (m *MockTaskService) UpdateTask(id int) error {
	return m.UpdateTaskFunc(id)
}

func (m *MockTaskService) DeleteTask(id int) error {
	return m.DeleteTaskFunc(id)
}

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name           string
		mockService    *MockTaskService
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			mockService: &MockTaskService{
				GetAllTasksFunc: func() ([]models.Task, error) {
					return []models.Task{
						{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
						{ID: 2, Title: "Task 2", UserID: 1, Completed: false},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"title":"Task 1","user_id":1,"completed":false},{"id":2,"title":"Task 2","user_id":1,"completed":false}]`,
		},
		{
			name: "internal server error",
			mockService: &MockTaskService{
				GetAllTasksFunc: func() ([]models.Task, error) {
					return nil, errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTaskHandler(tt.mockService)
			req := httptest.NewRequest("GET", "/tasks", nil)
			w := httptest.NewRecorder()

			handler.GetAll(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		mockService    *MockTaskService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			taskID: "1",
			mockService: &MockTaskService{
				GetTaskByIDFunc: func(id int) (models.Task, error) {
					return models.Task{ID: id, Title: "Task 1", UserID: 1}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Task 1","user_id":1,"completed":false}`,
		},
		{
			name:           "invalid id",
			taskID:         "abc",
			mockService:    &MockTaskService{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid task ID\n",
		},
		{
			name:   "not found",
			taskID: "999",
			mockService: &MockTaskService{
				GetTaskByIDFunc: func(id int) (models.Task, error) {
					return models.Task{}, errors.New("task not found")
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "task not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTaskHandler(tt.mockService)
			req := httptest.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			req.SetPathValue("id", tt.taskID)
			w := httptest.NewRecorder()

			handler.GetById(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockService    *MockTaskService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success",
			requestBody: `{"title":"New Task","user_id":1}`,
			mockService: &MockTaskService{
				CreateTaskFunc: func(task models.Task) error {
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:        "internal server error",
			requestBody: `{"title":"New Task","user_id":1}`,
			mockService: &MockTaskService{
				CreateTaskFunc: func(task models.Task) error {
					return errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTaskHandler(tt.mockService)
			req := httptest.NewRequest("POST", "/tasks", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			handler.CreateTask(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		mockService    *MockTaskService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			taskID: "1",
			mockService: &MockTaskService{
				DeleteTaskFunc: func(id int) error {
					return nil
				},
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "invalid id",
			taskID:         "abc",
			mockService:    &MockTaskService{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid task ID\n",
		},
		{
			name:   "internal server error",
			taskID: "1",
			mockService: &MockTaskService{
				DeleteTaskFunc: func(id int) error {
					return errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTaskHandler(tt.mockService)
			req := httptest.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			req.SetPathValue("id", tt.taskID)
			w := httptest.NewRecorder()

			handler.DeleteById(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMarkCompleteById(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		mockService    *MockTaskService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			taskID: "1",
			mockService: &MockTaskService{
				UpdateTaskFunc: func(id int) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "invalid id",
			taskID:         "abc",
			mockService:    &MockTaskService{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid task ID\n",
		},
		{
			name:   "internal server error",
			taskID: "1",
			mockService: &MockTaskService{
				UpdateTaskFunc: func(id int) error {
					return errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTaskHandler(tt.mockService)
			req := httptest.NewRequest("PATCH", "/tasks/"+tt.taskID+"/complete", nil)
			req.SetPathValue("id", tt.taskID)
			w := httptest.NewRecorder()

			handler.MarkCompleteById(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
