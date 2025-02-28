package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/emanuelquerty/gymulty/http/middleware"
)

type UserHandler struct {
	store domain.Store
	http.Handler
	logger *slog.Logger
}

func NewUserHandler(logger *slog.Logger, store domain.Store) *UserHandler {
	router := http.NewServeMux()
	userHandler := &UserHandler{
		store:   store,
		Handler: middleware.StripSlashes(router),
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
	e := &appError{Logger: u.logger}
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	user, err := u.store.GetUserByID(r.Context(), tenantID, userID)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	res := Response[[]domain.PublicUser]{
		Count: 1,
		Data: []domain.PublicUser{
			MapToPublicUser(user),
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) *appError {
	e := &appError{Logger: u.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	var user domain.User
	json.NewDecoder(r.Body).Decode(&user)

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	newUser, err := u.store.CreateUser(r.Context(), tenantID, user)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	res := Response[[]domain.PublicUser]{
		Count: 1,
		Data: []domain.PublicUser{
			MapToPublicUser(newUser),
		},
	}
	resourceURI := fmt.Sprintf("%s://%s%s/%d", r.URL.Scheme, r.Host, r.URL.String(), newUser.ID)

	w.Header().Set("Location", resourceURI)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) *appError {
	e := &appError{Logger: u.logger}
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	var update domain.UserUpdate
	json.NewDecoder(r.Body).Decode(&update)

	if update.Password != nil {
		hash, err := HashPassword(*update.Password)
		if err != nil {
			return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
		}
		update.Password = &hash
	}

	user, err := u.store.UpdateUser(r.Context(), tenantID, userID, update)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	res := Response[[]domain.PublicUser]{
		Count: 1,
		Data: []domain.PublicUser{
			MapToPublicUser(user),
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) deleteUserByID(w http.ResponseWriter, r *http.Request) *appError {
	e := &appError{Logger: u.logger}
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	err = u.store.DeleteUserByID(r.Context(), tenantID, userID)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	res := Response[any]{
		Count: 1,
		Data:  nil,
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(res)
	return nil
}

func (u *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) *appError {
	e := &appError{Logger: u.logger}
	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	users, err := u.store.GetAllUsers(r.Context(), tenantID)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	userCount := len(users)
	res := Response[[]domain.PublicUser]{
		Count: userCount,
		Data:  MapToPublicUsers(users),
	}

	w.WriteHeader(http.StatusOK)
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
	pubUsers := []domain.PublicUser{} // we want to return empty slice when len(users)==0, not nil slice
	for _, val := range users {
		curr := MapToPublicUser(val)
		pubUsers = append(pubUsers, curr)
	}
	return pubUsers
}
