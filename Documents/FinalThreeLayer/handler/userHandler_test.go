package handler

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

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		user         models.User
		setupMock    func(*MockuserService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success - create user",
			user: models.User{ID: 1, Name: "John Doe"},
			setupMock: func(mock *MockuserService) {
				mock.EXPECT().
					CreateUser(models.User{ID: 1, Name: "John Doe"}).
					Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "error - invalid payload",
			user: models.User{},
			setupMock: func(mock *MockuserService) {
				// No expectations for invalid payload
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "unexpected end of JSON input\n",
		},
		{
			name: "error - database failure",
			user: models.User{ID: 1, Name: "John Doe"},
			setupMock: func(mock *MockuserService) {
				mock.EXPECT().
					CreateUser(models.User{ID: 1, Name: "John Doe"}).
					Return(errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockuserService(ctrl)
			tt.setupMock(mockService)

			handler := NewUserHandler(mockService)
			body, _ := json.Marshal(tt.user)
			req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestUserHandler_GetByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		userID       string
		setupMock    func(*MockuserService)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "success - get user by ID",
			userID: "1",
			setupMock: func(mock *MockuserService) {
				mock.EXPECT().
					GetUserByID(1).
					Return(models.User{ID: 1, Name: "John Doe"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"John Doe"}`,
		},
		{
			name:   "error - invalid user ID",
			userID: "abc",
			setupMock: func(mock *MockuserService) {
				// No expectations for invalid ID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid user ID\n",
		},
		{
			name:   "error - user not found",
			userID: "999",
			setupMock: func(mock *MockuserService) {
				mock.EXPECT().
					GetUserByID(999).
					Return(models.User{}, errors.New("user not found"))
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockuserService(ctrl)
			tt.setupMock(mockService)

			handler := NewUserHandler(mockService)
			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)
			w := httptest.NewRecorder()

			handler.GetByUserId(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				if w.Code == http.StatusOK {
					assert.JSONEq(t, tt.expectedBody, w.Body.String())
				} else {
					assert.Equal(t, tt.expectedBody, w.Body.String())
				}
			}
		})
	}
}
