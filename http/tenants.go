package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/emanuelquerty/gymulty/http/middleware"
)

type TenantHandler struct {
	store domain.Store
	http.Handler
	logger *slog.Logger
}

func NewTenantHandler(logger *slog.Logger, store domain.Store) *TenantHandler {
	router := http.NewServeMux()
	handler := &TenantHandler{
		store:   store,
		Handler: middleware.StripSlashes(router),
		logger:  logger,
	}

	handler.registerRoutes(router)
	return handler
}

func (t *TenantHandler) registerRoutes(router *http.ServeMux) {
	router.Handle("POST /api/tenants/signup", errorHandler(t.createTenant))
}

func (t *TenantHandler) createTenant(w http.ResponseWriter, r *http.Request) *appError {
	var body domain.TenantRequestBody
	json.NewDecoder(r.Body).Decode(&body)

	tenant := domain.Tenant{
		BusinessName: body.BusinessName,
		Subdomain:    body.Subdomain,
	}
	newTenant, err := t.store.CreateTenant(r.Context(), tenant)
	e := &appError{Logger: t.logger}
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
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
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	newUser, err := t.store.CreateUser(r.Context(), newTenant.ID, userBody)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
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
