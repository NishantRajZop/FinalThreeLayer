package main

import (
	"log"
	"net/http"
	"user-service/db"
	"user-service/handlers"
)

func main() {
	db.InitDB("root:7280@tcp(localhost:3306)/userService")
	http.HandleFunc("/users", handlers.UsersHandler)
	http.HandleFunc("/users/", handlers.UserHandler)
	log.Println("User service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
