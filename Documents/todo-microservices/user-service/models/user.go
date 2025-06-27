package models

import "user-service/db"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateUser(user User) {
	db.DB.Exec("INSERT INTO users(name) VALUES(?)", user.Name)
}

func GetUserByID(id int) User {
	var user User
	db.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name)
	return user
}
