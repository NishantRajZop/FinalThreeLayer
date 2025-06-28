package store

import (
	"FinalThreeLayer/models"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type taskStore struct {
	db *sql.DB
}

// this new Function is like  a constructor  , if You will call this new with a db then it will return you a taskStore instance
// which will have all the CRUD functionalites🥳

// usually this is called by a service layer , so that service layer can have a kind of stick that can operate CRUD

func New(db *sql.DB) *taskStore {
	return &taskStore{db: db}
}

func (s *taskStore) GetAllTasks() ([]models.Task, error) {
	rows, _ := s.db.Query("SELECT id, title, user_id, completed FROM tasks")
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.Title, &t.UserID, &t.Completed)
		tasks = append(tasks, t)
	}
	fmt.Println(tasks)
	return tasks, nil
}

func (s *taskStore) GetTaskByID(id int) (models.Task, error) {
	var task models.Task
	s.db.QueryRow("SELECT id, title, user_id FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.UserID)
	return task, nil
}

func (s *taskStore) CreateTask(task models.Task) error {
	_, err := s.db.Exec("INSERT INTO tasks(title, user_id) VALUES(?, ?)", task.Title, task.UserID)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (s *taskStore) UpdateTask(id int) error {
	_, err := s.db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", true, id)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (s *taskStore) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}
