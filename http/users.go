package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/emanuelquerty/multency/domain"
)

type UserHandler struct {
	store domain.UserStore
	http.Handler
	logger *slog.Logger
}

func NewUserHandler(logger *slog.Logger, store domain.UserStore) *UserHandler {
	router := http.NewServeMux()
	userHandler := &UserHandler{
		store:   store,
		Handler: router,
		logger:  logger,
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
		return &appError{Error: err, Message: "invalid user id", Code: 400, Logger: u.logger}
	}

	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "invalid tenant id", Code: 400, Logger: u.logger}
	}

	user, err := u.store.GetUserByID(r.Context(), tenantID, userID)
	if err != nil {
		return &appError{Error: err, Message: "user was not found", Code: 404, Logger: u.logger}
	}

	res := Response[domain.PublicUser]{
		Success: true,
		Count:   1,
		Type:    "users",
		Data:    MapToPublicUser(user),
	}

	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")
	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "invalid tenant id", Code: 400, Logger: u.logger}
	}

	var user domain.User
	json.NewDecoder(r.Body).Decode(&user)

	err = user.HashPassword()
	if err != nil {
		return &appError{Error: err, Message: "could not create user", Code: 500, Logger: u.logger}
	}

	newUser, err := u.store.CreateUser(r.Context(), tenantID, user)
	if err != nil {
		return &appError{Error: err, Message: "could not create user", Code: 500, Logger: u.logger}
	}

	res := Response[domain.PublicUser]{
		Success: true,
		Count:   1,
		Type:    "users",
		Data:    MapToPublicUser(newUser),
	}
	resourceURI := fmt.Sprintf("%s://%s%s/%d", r.URL.Scheme, r.Host, r.URL.String(), newUser.ID)

	w.Header().Set("Location", resourceURI)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")
	uID := r.PathValue("userID")

	userID, err := strconv.Atoi(uID)
	if err != nil {
		return &appError{Error: err, Message: "invalid user id", Code: 400, Logger: u.logger}
	}

	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "invalid tenant id", Code: 400, Logger: u.logger}
	}

	var update domain.UserUpdate
	json.NewDecoder(r.Body).Decode(&update)

	user, err := u.store.UpdateUser(r.Context(), tenantID, userID, update)
	if err != nil {
		return &appError{Error: err, Message: "could not update user", Code: 404, Logger: u.logger}
	}

	res := Response[domain.PublicUser]{
		Success: true,
		Count:   1,
		Type:    "users",
		Data:    MapToPublicUser(user),
	}
	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) deleteUserByID(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")
	uID := r.PathValue("userID")

	userID, err := strconv.Atoi(uID)
	if err != nil {
		return &appError{Error: err, Message: "invalid user id", Code: 400, Logger: u.logger}
	}

	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "invalid tenant id", Code: 400, Logger: u.logger}
	}

	err = u.store.DeleteUserByID(r.Context(), tenantID, userID)
	if err != nil {
		return &appError{Error: err, Message: "could not delete user", Code: 404, Logger: u.logger}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (u *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) *appError {
	tID := r.PathValue("tenantID")

	tenantID, err := strconv.Atoi(tID)
	if err != nil {
		return &appError{Error: err, Message: "invalid tenant id", Code: 400, Logger: u.logger}
	}

	users, err := u.store.GetAllUsers(r.Context(), tenantID)
	if err != nil {
		return &appError{Error: err, Message: "could not retrieve users", Code: 404, Logger: u.logger}
	}

	userCount := len(users)
	res := Response[[]domain.PublicUser]{
		Success: true,
		Count:   userCount,
		Type:    "users",
		Data:    MapToPublicUsers(users),
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

func MapToPublicUsers(users []domain.User) []domain.PublicUser {
	var publicUsers []domain.PublicUser
	for _, val := range users {
		curr := MapToPublicUser(val)
		publicUsers = append(publicUsers, curr)
	}
	return publicUsers
}
