package userHandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"FinalThreeLayer/models"
)

type userHandler struct {
	service userService
}

func NewUserHandler(service userService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" && user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unexpected end of JSON input\n"))
		return
	}

	if err := h.service.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *userHandler) GetByUserId(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
