package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task-service/models"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tasks := models.GetAllTasks()
		json.NewEncoder(w).Encode(tasks)
	case "POST":
		var task models.Task
		json.NewDecoder(r.Body).Decode(&task)
		err := models.CreateTask(task)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Task Was Created Successfully"))
	}
}

func TaskHandlerId(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		task := models.GetTaskByID(id)
		json.NewEncoder(w).Encode(task)
	case "PUT":
		// var task models.Task
		//json.NewDecoder(r.Body).Decode(&task) // we will Not update Body Now , just true the Completed status
		err := models.UpdateTask(id)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Marked Comleted The Task!!"))
		}
	case "DELETE":
		err := models.DeleteTask(id)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Task Has been Deleted!!"))
		}
	}
}
