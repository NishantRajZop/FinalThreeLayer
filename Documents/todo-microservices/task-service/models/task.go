package models

import (
	"fmt"
	"task-service/db"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	UserID    int    `json:"user_id"`
	Completed bool   `json:"completed"`
}

func CreateTask(task Task) error { // for Creating a Task U will need a TaskTitle and User id
	_, err := db.DB.Exec("INSERT INTO tasks(title, user_id) VALUES(?, ?)", task.Title, task.UserID)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func GetAllTasks() []Task {
	rows, _ := db.DB.Query("SELECT id, title, user_id, completed FROM tasks")
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.UserID, &t.Completed)
		tasks = append(tasks, t)
	}
	fmt.Println(tasks)
	return tasks
}

func GetTaskByID(id int) Task {
	var task Task
	db.DB.QueryRow("SELECT id, title, user_id FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.UserID)
	return task
}

func UpdateTask(id int) error {
	_, err := db.DB.Exec("UPDATE tasks SET completed = ? WHERE id = ?", true, id)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func DeleteTask(id int) error {
	_, err := db.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}
