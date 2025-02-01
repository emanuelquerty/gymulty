package http

import (
	"encoding/json"
	"fmt"
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
	router.Handle("GET /tenants/{tenantID}/users/{userID}", errorHandler(u.getUserByID))

	// router.Handle("GET users/{id}", errorHandler(u.getUserByID))
	router.Handle("PUT /users/{id}", errorHandler(u.updateUser))
	router.Handle("POST /users", errorHandler(u.createUser))
	router.Handle("DELETE /users/{id}", errorHandler(u.deleteUserByID))
}

func (u *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")
	uID := r.PathValue("userID")

	userID, err := strconv.Atoi(uID)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid user id", Code: http.StatusBadRequest}
	}

	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid tenant id", Code: http.StatusBadRequest}
	}

	user, err := u.store.GetUserByID(tenantID, userID)
	if err != nil {
		return &appError{Error: err, Message: "user was not found", Code: http.StatusNotFound}
	}

	json.NewEncoder(w).Encode(user)
	return nil
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) *appError {
	var user domain.User
	json.NewDecoder(r.Body).Decode(&user)

	err := user.HashPassword()
	if err != nil {
		return &appError{Error: err, Message: "could not create user", Code: 500}
	}

	newUser, err := u.store.CreateUser(user)
	if err != nil {
		return &appError{Error: err, Message: "could not create user", Code: 500}
	}

	publicUser := MapToPublicUser(newUser)
	resourceURI := fmt.Sprintf("%s://%s%s/%d", r.URL.Scheme, r.Host, r.URL.String(), publicUser.ID)

	w.Header().Set("Location", resourceURI)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(publicUser)
	return nil
}

func (u *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) *appError {
	idString := r.PathValue("id")

	userID, err := strconv.Atoi(idString)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid id", Code: http.StatusBadRequest}
	}

	var update domain.UserUpdate
	json.NewDecoder(r.Body).Decode(&update)

	updatedUser, err := u.store.UpdateUser(userID, update)
	if err != nil {
		return &appError{Error: err, Message: "could not update user", Code: http.StatusNotFound}
	}
	json.NewEncoder(w).Encode(updatedUser)
	return nil
}

func (u *UserHandler) deleteUserByID(w http.ResponseWriter, r *http.Request) *appError {
	idString := r.PathValue("id")

	userID, err := strconv.Atoi(idString)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid id", Code: http.StatusBadRequest}
	}

	err = u.store.DeleteUserByID(userID)
	if err != nil {
		return &appError{Error: err, Message: "could not delete user", Code: http.StatusNotFound}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func MapToPublicUser(user domain.User) domain.PublicUser {
	return domain.PublicUser{
		ID:        user.ID,
		TenantID:  user.TenantID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}
}
