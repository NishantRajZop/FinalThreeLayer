package service

import (
	"FinalThreeLayer/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
		{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
	}

	tests := []struct {
		name        string
		setupMock   func(*MockTaskStore)
		wantTasks   []models.Task
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - get all tasks",
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().GetAllTasks().Return(expected, nil)
			},
			wantTasks: expected,
			wantErr:   false,
		},
		{
			name: "error - database failure",
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().GetAllTasks().Return(nil, errors.New("database error"))
			},
			wantTasks:   nil,
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockStore := NewMockTaskStore(ctrl)
			tt.setupMock(mockStore)

			service := NewTaskService(mockStore)
			tasks, err := service.GetAllTasks()

			if err == nil && !reflect.DeepEqual(tasks, expected) {
				t.Errorf("expected %v , got %v", expected, tasks)
			}

			if err != nil && !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("expected %v , got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		taskID      int
		setupMock   func(*MockTaskStore)
		wantTask    models.Task
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - get task by ID",
			taskID: 1,
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().GetTaskByID(1).Return(models.Task{ID: 1, Title: "Task 1", UserID: 1}, nil)
			},
			wantTask: models.Task{ID: 1, Title: "Task 1", UserID: 1},
			wantErr:  false,
		},
		{
			name:   "error - task not found",
			taskID: 999,
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().GetTaskByID(999).Return(models.Task{}, errors.New("task not found"))
			},
			wantTask:    models.Task{},
			wantErr:     true,
			expectedErr: errors.New("task not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockTaskStore(ctrl)
			tt.setupMock(mockStore)

			service := NewTaskService(mockStore)
			task, err := service.GetTaskByID(tt.taskID)

			if tt.wantErr { // if there is a error then Compare that Error types and errorMessages are same
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else { // if there is no error means you got a res body in return , then do taskBody is same as expected or not
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTask, task)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		task        models.Task
		setupMock   func(*MockTaskStore)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - create task",
			task: models.Task{Title: "New Task", UserID: 1},
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().CreateTask(models.Task{Title: "New Task", UserID: 1}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - database failure",
			task: models.Task{Title: "New Task", UserID: 1},
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().CreateTask(models.Task{Title: "New Task", UserID: 1}).Return(errors.New("database error"))
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name: "error - empty title",
			task: models.Task{Title: "", UserID: 1},
			setupMock: func(mock *MockTaskStore) { // this should be Handle at validation in service layer
			},
			wantErr:     true,
			expectedErr: errors.New("task title cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockTaskStore(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStore)
			}

			service := NewTaskService(mockStore)
			err := service.CreateTask(tt.task)

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

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		taskID      int
		setupMock   func(*MockTaskStore)
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - update task",
			taskID: 1,
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().UpdateTask(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - database failure",
			taskID: 1,
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().UpdateTask(1).Return(errors.New("database error"))
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name:   "error - invalid task ID",
			taskID: 0,
			setupMock: func(mock *MockTaskStore) {
				// No expectation as validation should fail before calling the store
			},
			wantErr:     true,
			expectedErr: errors.New("invalid task ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockTaskStore(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStore)
			}

			service := NewTaskService(mockStore)
			err := service.UpdateTask(tt.taskID)

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

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		taskID      int
		setupMock   func(*MockTaskStore)
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - delete task",
			taskID: 1,
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().DeleteTask(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - database failure",
			taskID: 1,
			setupMock: func(mock *MockTaskStore) {
				mock.EXPECT().DeleteTask(1).Return(errors.New("database error"))
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name:   "error - invalid task ID",
			taskID: 0,
			setupMock: func(mock *MockTaskStore) {
				// No expectation as validation should fail before calling the store
			},
			wantErr:     true,
			expectedErr: errors.New("invalid task ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockTaskStore(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStore)
			}

			service := NewTaskService(mockStore)
			err := service.DeleteTask(tt.taskID)

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
