package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/multency/domain"
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
	if err != nil {
		return &appError{Error: err, Message: "could not create tenant", Code: 400, Logger: t.logger}
	}

	userBody := domain.User{
		TenantID:  newTenant.ID,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Role:      "admin",
	}
	err = userBody.HashPassword()
	if err != nil {
		return &appError{Error: err, Message: "could not create tenant", Code: 500, Logger: t.logger}
	}

	newUser, err := t.userStore.CreateUser(r.Context(), newTenant.ID, userBody)
	if err != nil {
		return &appError{Error: err, Message: "could not create user for given tenant", Code: 400, Logger: t.logger}
	}

	res := Response[TenantSignupResponse]{
		Success: true,
		Count:   1,
		Type:    "tenants",
		Data: TenantSignupResponse{
			Message: "tenant registered successfully",
			Tenant:  newTenant,
			Admin:   MapToPublicUser(newUser),
		},
	}
	json.NewEncoder(w).Encode(res)
	return nil
}
