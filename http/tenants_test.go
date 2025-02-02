package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emanuelquerty/multency/domain"
	"github.com/emanuelquerty/multency/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTenant(t *testing.T) {
	body := domain.TenantBody{
		BusinessName: "SwoleGym",
		Subdomain:    "swolegym",
		FirstName:    "Peter",
		LastName:     "Gray",
		Email:        "pgray@email.com",
		Password:     "ReallySecret1001",
	}

	jsonBody, _ := json.Marshal(body)
	buf := bytes.NewBuffer(jsonBody)

	t.Run("returns newly created tenant on success", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(ctx context.Context, tenantID int, data domain.User) (domain.User, error) {
			data.ID = 1
			return data, nil
		}

		tenantStore := new(mock.TenantStore)
		tenantStore.CreateTenantFn = func(ctx context.Context, data domain.Tenant) (domain.Tenant, error) {
			data.ID = 1
			return data, nil
		}

		req := httptest.NewRequest("POST", "/api/tenants/signup", buf)
		res := newTenantRequest(tenantStore, userStore, req)

		want := TenantSignupResponse{
			Message: "tenant registered successfully",
			Tenant: domain.Tenant{
				ID:           1,
				BusinessName: body.BusinessName,
				Subdomain:    body.Subdomain,
			},
			Admin: domain.PublicUser{
				ID:        1,
				TenantID:  1,
				FirstName: body.FirstName,
				LastName:  body.LastName,
				Role:      "admin",
			},
		}

		var got TenantSignupResponse
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "tenants should match")
	})
}

func newTenantRequest(tenantStore *mock.TenantStore, userstore *mock.UserStore, req *http.Request) *httptest.ResponseRecorder {
	handler := NewTenantHandler(tenantStore, userstore)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res
}
