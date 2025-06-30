package store

import (
	"FinalThreeLayer/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	t.Helper()

	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
}

func TestCreateUser(t *testing.T) {
	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	testUser := models.User{
		ID:   1,
		Name: "test_user1 1",
	}

	// Creates a new instance of your repository/struct (presumably named str)
	// Passes the mock database connection to it
	mock.ExpectExec("INSERT INTO users(name) VALUES(?)").WithArgs(testUser.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.CreateUser(testUser)

	if err != nil {
		t.Errorf("expected success message, got: %s", err)
	}
}

func TestGetUserByID(t *testing.T) {

	db, mock, err := getMockDB(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := New(db)

	testID := 1
	expectedUser := models.User{
		ID:   testID,
		Name: "test_user1",
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedUser.ID, expectedUser.Name)

	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(testID).
		WillReturnRows(rows)

	user, err := s.GetUserByID(testID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf("got user ID %d, expected %d", user.ID, expectedUser.ID)
	}
	if user.Name != expectedUser.Name {
		t.Errorf("got user name '%s', expected '%s'", user.Name, expectedUser.Name)
	}

	//if err := mock.ExpectationsWereMet(); err != nil {
	//	t.Errorf("there were unfulfilled expectations: %s", err)
	//}
}
