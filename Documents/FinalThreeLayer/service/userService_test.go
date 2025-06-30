package service

import (
	"FinalThreeLayer/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		user        models.User
		setupMock   func(*MockUserStore)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - create user",
			user: models.User{ID: 1, Name: "John Doe"},
			setupMock: func(mock *MockUserStore) {
				mock.EXPECT().
					CreateUser(models.User{ID: 1, Name: "John Doe"}).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "error - database failure",
			user: models.User{ID: 1, Name: "John Doe"},
			setupMock: func(mock *MockUserStore) {
				mock.EXPECT().
					CreateUser(models.User{ID: 1, Name: "John Doe"}).
					Return(errors.New("database error")).
					Times(1)
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name: "error - empty name",
			user: models.User{ID: 1, Name: ""},
			setupMock: func(mock *MockUserStore) {
				// No expectation as validation should fail before calling the store
			},
			wantErr:     true,
			expectedErr: errors.New("user name cannot be empty"),
		},
		{
			name: "error - invalid ID",
			user: models.User{ID: 0, Name: "John Doe"},
			setupMock: func(mock *MockUserStore) {
				// No expectation as validation should fail before calling the store
			},
			wantErr:     true,
			expectedErr: errors.New("invalid user ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockUserStore(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStore)
			}

			service := NewUserService(mockStore)
			err := service.CreateUser(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		userID      int
		setupMock   func(*MockUserStore)
		wantUser    models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - get user by ID",
			userID: 1,
			setupMock: func(mock *MockUserStore) {
				mock.EXPECT().
					GetUserByID(1).
					Return(models.User{ID: 1, Name: "John Doe"}, nil).
					Times(1)
			},
			wantUser: models.User{ID: 1, Name: "John Doe"},
			wantErr:  false,
		},
		{
			name:   "error - user not found",
			userID: 999,
			setupMock: func(mock *MockUserStore) {
				mock.EXPECT().
					GetUserByID(999).
					Return(models.User{}, errors.New("user not found")).
					Times(1)
			},
			wantUser:    models.User{},
			wantErr:     true,
			expectedErr: errors.New("user not found"),
		},
		{
			name:   "error - invalid ID",
			userID: 0,
			setupMock: func(mock *MockUserStore) {
				// No expectation as validation should fail before calling the store
			},
			wantUser:    models.User{},
			wantErr:     true,
			expectedErr: errors.New("invalid user ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockUserStore(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStore)
			}

			service := NewUserService(mockStore)
			user, err := service.GetUserByID(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}
		})
	}
}
