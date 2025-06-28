package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	UserID    int    `json:"user_id"`
	Completed bool   `json:"completed"`
}
