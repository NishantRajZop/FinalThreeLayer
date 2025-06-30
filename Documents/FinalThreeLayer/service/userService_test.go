package service

import (
	"FinalThreeLayer/models"
	"errors"
	"testing"
)

// MockUserStore implements UserStore for testing
type MockUserStore struct {
	CreateUserFunc  func(models.User) error
	GetUserByIDFunc func(id int) (models.User, error)
}

func (m *MockUserStore) CreateUser(user models.User) error {
	return m.CreateUserFunc(user)
}

func (m *MockUserStore) GetUserByID(id int) (models.User, error) {
	return m.GetUserByIDFunc(id)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        models.User
		mockStore   *MockUserStore
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - create user",
			user: models.User{ID: 1, Name: "John Doe"},
			mockStore: &MockUserStore{
				CreateUserFunc: func(user models.User) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "error - database failure",
			user: models.User{ID: 1, Name: "John Doe"},
			mockStore: &MockUserStore{
				CreateUserFunc: func(user models.User) error {
					return errors.New("database error")
				},
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		//{
		//	name: "error - empty name",
		//	user: models.User{ID: 1, Name: ""},
		//	mockStore: &MockUserStore{
		//		CreateUserFunc: func(user models.User) error {
		//			return nil // Shouldn't be called
		//		},
		//	},
		//	wantErr:     true,
		//	expectedErr: errors.New("user name cannot be empty"),
		//},
		//{
		//	name: "error - invalid ID",
		//	user: models.User{ID: 0, Name: "John Doe"},
		//	mockStore: &MockUserStore{
		//		CreateUserFunc: func(user models.User) error {
		//			return nil // Shouldn't be called
		//		},
		//	},
		//	wantErr:     true,
		//	expectedErr: errors.New("invalid user ID"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(tt.mockStore)
			err := service.CreateUser(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("CreateUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		mockStore   *MockUserStore
		wantUser    models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - get user by ID",
			userID: 1,
			mockStore: &MockUserStore{
				GetUserByIDFunc: func(id int) (models.User, error) {
					return models.User{ID: id, Name: "John Doe"}, nil
				},
			},
			wantUser: models.User{ID: 1, Name: "John Doe"},
			wantErr:  false,
		},
		{
			name:   "error - user not found",
			userID: 999,
			mockStore: &MockUserStore{
				GetUserByIDFunc: func(id int) (models.User, error) {
					return models.User{}, errors.New("user not found")
				},
			},
			wantUser:    models.User{},
			wantErr:     true,
			expectedErr: errors.New("user not found"),
		},
		//{
		//	name:   "error - invalid ID",
		//	userID: 0,
		//	mockStore: &MockUserStore{
		//		GetUserByIDFunc: func(id int) (models.User, error) {
		//			return models.User{}, nil // Shouldn't be called
		//		},
		//	},
		//	wantUser:    models.User{},
		//	wantErr:     true,
		//	expectedErr: errors.New("invalid user ID"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(tt.mockStore)
			user, err := service.GetUserByID(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("GetUserByID() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if user != tt.wantUser {
				t.Errorf("GetUserByID() = %v, want %v", user, tt.wantUser)
			}
		})
	}
}
