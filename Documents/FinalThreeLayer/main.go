package main

import (
	"FinalThreeLayer/dataSource"
	taskHandler2 "FinalThreeLayer/handler/taskHandler"
	"FinalThreeLayer/handler/userHandler"
	taskStore "FinalThreeLayer/store/task"
	userStore "FinalThreeLayer/store/user"

	taskService "FinalThreeLayer/service"
	userService "FinalThreeLayer/service"

	"log"
	"net/http"
)

func main() {
	db := dataSource.InitDB("root:7280@tcp(localhost:3306)/taskService")
	tStore := taskStore.New(db)
	uStore := userStore.New(db)

	tService := taskService.NewTaskService(tStore)
	uService := userService.NewUserService(uStore)

	tHandler := taskHandler2.NewTaskHandler(tService)
	uHandler := userHandler.NewUserHandler(uService)

	http.HandleFunc("GET /tasks", tHandler.GetAll)
	http.HandleFunc("GET /tasks/{id}", tHandler.GetById)
	http.HandleFunc("POST /tasks", tHandler.CreateTask)
	http.HandleFunc("DELETE /tasks/{id}", tHandler.DeleteById)
	http.HandleFunc("PUT /tasks/{id}", tHandler.MarkCompleteById)
	http.HandleFunc("POST /users", uHandler.CreateUser)
	http.HandleFunc("GET /users/{id}", uHandler.GetByUserId)

	log.Println("services running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
