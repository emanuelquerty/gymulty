package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/gymulty/domain"
)

type TenantHandler struct {
	tenantStore domain.TenantStore
	userStore   domain.UserStore
	http.Handler
	logger *slog.Logger
}

func NewTenantHandler(logger *slog.Logger, tenantStore domain.TenantStore, userStore domain.UserStore) *TenantHandler {
	router := http.NewServeMux()
	handler := &TenantHandler{
		tenantStore: tenantStore,
		userStore:   userStore,
		Handler:     router,
		logger:      logger,
	}

	handler.registerRoutes(router)
	return handler
}

func (t *TenantHandler) registerRoutes(router *http.ServeMux) {
	userHandler := NewUserHandler(t.logger, t.userStore)

	router.Handle("POST /api/tenants/signup", errorHandler(t.createTenant))
	router.Handle("/api/tenants/{tenantID}/", userHandler)
}

func (t *TenantHandler) createTenant(w http.ResponseWriter, r *http.Request) *appError {
	var body domain.TenantRequestBody
	json.NewDecoder(r.Body).Decode(&body)

	tenant := domain.Tenant{
		BusinessName: body.BusinessName,
		Subdomain:    body.Subdomain,
	}
	newTenant, err := t.tenantStore.CreateTenant(r.Context(), tenant)
	e := &appError{Logger: t.logger}
	if err != nil {
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrInternal)
	}

	userBody := domain.User{
		TenantID:  newTenant.ID,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Role:      "admin",
	}
	userBody.Password, err = HashPassword(userBody.Password)
	if err != nil {
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrInternal)
	}

	newUser, err := t.userStore.CreateUser(r.Context(), newTenant.ID, userBody)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.withContext(err, "A tenant with specified id was not found", ErrNotFound)
		}
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrInternal)
	}

	res := Response[TenantSignupResponse]{
		Count: 1,
		Data: TenantSignupResponse{
			Message: "tenant registered successfully",
			Tenant:  newTenant,
			Admin:   MapToPublicUser(newUser),
		},
	}
	json.NewEncoder(w).Encode(res)
	return nil
}
