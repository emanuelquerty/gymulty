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

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.Handler = http.StripPrefix("/api", u.Handler)
	u.Handler.ServeHTTP(w, r)
}

func (u *UserHandler) registerRoutes(router *http.ServeMux) {
	router.Handle("GET /users/{id}", errorHandler(u.getUserByID))
	router.Handle("POST /users", errorHandler(u.createUser))
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

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) *appError {
	var user domain.User
	json.NewDecoder(r.Body).Decode(&user)

	newUser, err := u.store.CreateUser(user)
	if err != nil {
		return &appError{Error: err, Message: "could create user", Code: 400}
	}

	publicUser := MapToPublicUser(newUser)

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(publicUser)
	return nil
}

func MapToPublicUser(user domain.User) domain.PublicUser {
	return domain.PublicUser{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Role:      user.Role,
	}
}
