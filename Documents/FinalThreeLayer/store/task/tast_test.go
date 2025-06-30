package store

import (
	"FinalThreeLayer/models"
	"database/sql"
	_ "database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	t.Helper()

	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
}

func TestGetAllTasks(t *testing.T) {
	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	expectedTasks := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
		{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "completed"}).
		AddRow(expectedTasks[0].ID, expectedTasks[0].Title, expectedTasks[0].UserID, expectedTasks[0].Completed).
		AddRow(expectedTasks[1].ID, expectedTasks[1].Title, expectedTasks[1].UserID, expectedTasks[1].Completed)

	mock.ExpectQuery("SELECT id, title, user_id, completed FROM tasks").WillReturnRows(rows)

	tasks, err := s.GetAllTasks()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(tasks) != len(expectedTasks) {
		t.Errorf("expected %d tasks, got %d", len(expectedTasks), len(tasks))
	}

	for i, task := range tasks {
		if task != expectedTasks[i] {
			t.Errorf("task %d mismatch: got %+v, expected %+v", i, task, expectedTasks[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTaskByID(t *testing.T) {
	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	testID := 1
	expectedTask := models.Task{
		ID:     testID,
		Title:  "Test Task",
		UserID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id"}).
		AddRow(expectedTask.ID, expectedTask.Title, expectedTask.UserID)

	mock.ExpectQuery("SELECT id, title, user_id FROM tasks WHERE id = ?").
		WithArgs(testID).
		WillReturnRows(rows)

	task, err := s.GetTaskByID(testID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if task != expectedTask {
		t.Errorf("task mismatch: got %+v, expected %+v", task, expectedTask)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateTask(t *testing.T) {
	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	testTask := models.Task{
		Title:  "New Task",
		UserID: 1,
	}

	mock.ExpectExec("INSERT INTO tasks(title, user_id) VALUES(?, ?)").
		WithArgs(testTask.Title, testTask.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.CreateTask(testTask)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	testID := 1

	mock.ExpectExec(`UPDATE tasks SET completed = ? WHERE id = ?`).
		WithArgs(true, testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = s.UpdateTask(testID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	testID := 1

	mock.ExpectExec("DELETE FROM tasks WHERE id = ?").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = s.DeleteTask(testID)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

}
