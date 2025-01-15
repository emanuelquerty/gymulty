package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emanuelquerty/multency/domain"
)

type UserHandler struct {
	store domain.UserStore
	http.Handler
}

func NewUserHandler(store domain.UserStore) *UserHandler {
	router := http.NewServeMux()
	userHandler := &UserHandler{
		store:   store,
		Handler: router,
	}
	userHandler.registerRoutes(router)
	return userHandler
}

func (u *UserHandler) registerRoutes(router *http.ServeMux) {
	router.Handle("GET /api/users/{id}", errorHandler(u.getUserByID))
}

func (u *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) *appError {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid id", Code: http.StatusBadRequest}
	}

	user, err := u.store.GetUserByID(id)
	if err != nil {
		return &appError{Error: err, Message: "user was not found", Code: http.StatusNotFound}
	}

	json.NewEncoder(w).Encode(user)
	return nil
}
