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

func (u *UserHandler) registerRoutes(router *http.ServeMux) {
	router.Handle("GET /api/tenants/{tenantID}/users/{userID}", errorHandler(u.getUserByID))
	router.Handle("POST /api/tenants/{tenantID}/users", errorHandler(u.createUser))

	router.Handle("PUT /api/tenants/{tenantID}/users/{userID}", errorHandler(u.updateUser))
	router.Handle("DELETE /api/tenants/{tenantID}/users/{userID}", errorHandler(u.deleteUserByID))
	router.Handle("GET /api/tenants/{tenantID}/users", errorHandler(u.getAllUsers))
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

	user, err := u.store.GetUserByID(r.Context(), tenantID, userID)
	if err != nil {
		return &appError{Error: err, Message: "user was not found", Code: http.StatusNotFound}
	}

	json.NewEncoder(w).Encode(user)
	return nil
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")
	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid tenant id", Code: http.StatusBadRequest}
	}

	var user domain.User
	json.NewDecoder(r.Body).Decode(&user)

	err = user.HashPassword()
	if err != nil {
		return &appError{Error: err, Message: "could not create user", Code: 500}
	}

	newUser, err := u.store.CreateUser(r.Context(), tenantID, user)
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

	var update domain.UserUpdate
	json.NewDecoder(r.Body).Decode(&update)

	updatedUser, err := u.store.UpdateUser(r.Context(), tenantID, userID, update)
	if err != nil {
		return &appError{Error: err, Message: "could not update user", Code: http.StatusNotFound}
	}
	json.NewEncoder(w).Encode(updatedUser)
	return nil
}

func (u *UserHandler) deleteUserByID(w http.ResponseWriter, r *http.Request) *appError {
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

	err = u.store.DeleteUserByID(r.Context(), tenantID, userID)
	if err != nil {
		return &appError{Error: err, Message: "could not delete user", Code: http.StatusNotFound}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (u *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")

	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "malformed url: invalid tenant id", Code: http.StatusBadRequest}
	}

	users, err := u.store.GetAllUsers(r.Context(), tenantID)
	if err != nil {
		return &appError{Error: err, Message: "could not retrieve users", Code: http.StatusNotFound}
	}

	userCount := len(users)
	res := struct {
		Message string
		Users   []domain.User
	}{
		Message: fmt.Sprintf("found %d users", userCount),
		Users:   users,
	}

	json.NewEncoder(w).Encode(res)
	return nil
}

func MapToPublicUser(user domain.User) domain.PublicUser {
	return domain.PublicUser{
		ID:        user.ID,
		TenantID:  user.TenantID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
