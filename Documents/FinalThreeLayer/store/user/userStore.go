package store

import (
	"FinalThreeLayer/models"
	"database/sql"
	"fmt"
)

type userStore struct {
	db *sql.DB
}

func New(db *sql.DB) *userStore {
	return &userStore{db: db}
}

func (s *userStore) CreateUser(user models.User) error {
	_, err := s.db.Exec("INSERT INTO users(name) VALUES(?)", user.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *userStore) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name)
	if err != nil {
		return models.User{}, fmt.Errorf("database error: %v", err)
	}

	return user, nil
}
