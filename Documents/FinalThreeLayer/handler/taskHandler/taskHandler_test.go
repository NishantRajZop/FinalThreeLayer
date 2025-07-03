package taskHandler

import (
	"FinalThreeLayer/models"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
		{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
	}
	tests := []struct {
		name         string
		setupMock    func(*MocktaskService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success - get all tasks",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().GetAllTasks().Return(expected, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"completed":false, "id":1, "title":"Task 1", "user_id":1},{"completed":true, "id":2, "title":"Task 2", "user_id":1}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMocktaskService(ctrl)
			tt.setupMock(mockService)

			handler := NewTaskHandler(mockService)
			req := httptest.NewRequest("GET", "/tasks", nil)
			w := httptest.NewRecorder()

			handler.GetAll(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestTaskHandler_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := models.Task{ID: 1, Title: "Task 1", UserID: 1, Completed: false}

	tests := []struct {
		name         string
		taskID       string
		setupMock    func(*MocktaskService)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "success - get task by ID",
			taskID: "1",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().GetTaskByID(1).Return(expected, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"completed":false, "id":1, "title":"Task 1", "user_id":1}`,
		},
		{
			name:         "error - invalid ID",
			taskID:       "abc",
			setupMock:    func(mock *MocktaskService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "{}",
		},
		{
			name:   "error - task not found",
			taskID: "999",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().GetTaskByID(999).Return(models.Task{}, errors.New("not found"))
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMocktaskService(ctrl)
			tt.setupMock(mockService)

			handler := NewTaskHandler(mockService)
			req := httptest.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			req.SetPathValue("id", tt.taskID)
			w := httptest.NewRecorder()

			handler.GetById(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			//if tt.expectedBody != "" {
			//	assert.JSONEq(t, tt.expectedBody, w.Body.String())
			//} else {
			//	assert.Equal(t, tt.expectedBody, w.Body.String())
			//}
		})
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		task         models.Task
		setupMock    func(*MocktaskService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success - create task",
			task: models.Task{Title: "New Task"},
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().CreateTask(models.Task{Title: "New Task"}).Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "error - invalid payload",
			task:         models.Task{},
			setupMock:    func(mock *MocktaskService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error - service failure",
			task: models.Task{Title: "New Task"},
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().CreateTask(models.Task{Title: "New Task"}).Return(errors.New("create failed"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "create failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMocktaskService(ctrl)
			tt.setupMock(mockService)

			handler := NewTaskHandler(mockService)
			body, _ := json.Marshal(tt.task)
			req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateTask(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestTaskHandler_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		taskID       string
		setupMock    func(*MocktaskService)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "success - delete task",
			taskID: "1",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().DeleteTask(1).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "error - invalid ID",
			taskID:       "abc",
			setupMock:    func(mock *MocktaskService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid task ID\n",
		},
		{
			name:   "error - service failure",
			taskID: "1",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().DeleteTask(1).Return(errors.New("delete failed"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "delete failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMocktaskService(ctrl)
			tt.setupMock(mockService)

			handler := NewTaskHandler(mockService)
			req := httptest.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			req.SetPathValue("id", tt.taskID)
			w := httptest.NewRecorder()

			handler.DeleteById(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestTaskHandler_MarkCompleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		taskID       string
		setupMock    func(*MocktaskService)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "success - mark task complete",
			taskID: "1",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().UpdateTask(1).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "error - invalid ID",
			taskID:       "abc",
			setupMock:    func(mock *MocktaskService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid task ID\n",
		},
		{
			name:   "error - service failure",
			taskID: "1",
			setupMock: func(mock *MocktaskService) {
				mock.EXPECT().UpdateTask(1).Return(errors.New("update failed"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "update failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMocktaskService(ctrl)
			tt.setupMock(mockService)

			handler := NewTaskHandler(mockService)
			req := httptest.NewRequest("PATCH", "/tasks/"+tt.taskID+"/complete", nil)
			req.SetPathValue("id", tt.taskID)
			w := httptest.NewRecorder()

			handler.MarkCompleteById(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
