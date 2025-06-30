package handler

import (
	"FinalThreeLayer/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockUserService implements userService for testing
type MockUserService struct {
	CreateUserFunc  func(models.User) error
	GetUserByIDFunc func(id int) (models.User, error)
}

func (m *MockUserService) CreateUser(user models.User) error {
	return m.CreateUserFunc(user)
}

func (m *MockUserService) GetUserByID(id int) (models.User, error) {
	return m.GetUserByIDFunc(id)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockService    *MockUserService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success",
			requestBody: `{"id":1,"name":"John Doe"}`,
			mockService: &MockUserService{
				CreateUserFunc: func(user models.User) error {
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:        "internal server error",
			requestBody: `{"id":1,"name":"John Doe"}`,
			mockService: &MockUserService{
				CreateUserFunc: func(user models.User) error {
					return errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService)
			req := httptest.NewRequest("POST", "/users", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetByUserId(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockService    *MockUserService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			userID: "1",
			mockService: &MockUserService{
				GetUserByIDFunc: func(id int) (models.User, error) {
					return models.User{ID: id, Name: "John Doe"}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"John Doe"}`,
		},
		{
			name:           "invalid id",
			userID:         "abc",
			mockService:    &MockUserService{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid user ID\n",
		},
		{
			name:   "not found",
			userID: "999",
			mockService: &MockUserService{
				GetUserByIDFunc: func(id int) (models.User, error) {
					return models.User{}, errors.New("user not found")
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService)
			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)
			w := httptest.NewRecorder()

			handler.GetByUserId(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
