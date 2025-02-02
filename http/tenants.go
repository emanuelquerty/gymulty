package http

import (
	"encoding/json"
	"net/http"

	"github.com/emanuelquerty/multency/domain"
)

type TenantHandler struct {
	tenantStore domain.TenantStore
	userStore   domain.UserStore
	http.Handler
}

func NewTenantHandler(tenantStore domain.TenantStore, userStore domain.UserStore) *TenantHandler {
	router := http.NewServeMux()
	handler := &TenantHandler{
		tenantStore: tenantStore,
		userStore:   userStore,
		Handler:     router,
	}

	handler.registerRoutes(router)
	return handler
}

func (t *TenantHandler) registerRoutes(router *http.ServeMux) {
	userHandler := NewUserHandler(t.userStore)

	router.Handle("POST /api/tenants/signup", errorHandler(t.createTenant))
	router.Handle("/api/tenants/{tenantID}/", userHandler)
}

func (t *TenantHandler) createTenant(w http.ResponseWriter, r *http.Request) *appError {
	var body domain.TenantBody
	json.NewDecoder(r.Body).Decode(&body)

	tenant := domain.Tenant{
		BusinessName: body.BusinessName,
		Subdomain:    body.Subdomain,
	}
	newTenant, err := t.tenantStore.CreateTenant(r.Context(), tenant)
	if err != nil {
		return &appError{Error: err, Message: "could not create tenant", Code: http.StatusBadRequest}
	}

	user := domain.User{
		TenantID:  newTenant.ID,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Role:      "admin",
	}
	err = user.HashPassword()
	if err != nil {
		return &appError{Error: err, Message: "could not create tenant", Code: http.StatusInternalServerError}
	}

	newUser, err := t.userStore.CreateUser(r.Context(), newTenant.ID, user)
	if err != nil {
		return &appError{Error: err, Message: "could not create user for given tenant", Code: http.StatusBadRequest}
	}

	res := TenantSignupResponse{
		Message: "tenant registered successfully",
		Tenant:  newTenant,
		Admin:   MapToPublicUser(newUser),
	}
	json.NewEncoder(w).Encode(res)
	return nil
}
