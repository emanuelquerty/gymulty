package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/emanuelquerty/gymulty/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTenant(t *testing.T) {
	body := domain.TenantRequestBody{
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
		store := new(mock.Store)
		store.CreateTenantFn = func(ctx context.Context, data domain.Tenant) (domain.Tenant, error) {
			data.ID = 1
			return data, nil
		}

		store.CreateUserFn = func(ctx context.Context, tenantID int, data domain.User) (domain.User, error) {
			data.ID = 1
			return data, nil
		}

		req := httptest.NewRequest("POST", "/api/tenants/signup", buf)
		res := newTenantRequest(store, req)

		want := Response[TenantSignupResponse]{
			Count: 1,
			Data: TenantSignupResponse{
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
			},
		}

		var got Response[TenantSignupResponse]
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "tenants should match")
	})
}

func newTenantRequest(store *mock.Store, req *http.Request) *httptest.ResponseRecorder {
	handler := NewTenantHandler(slog.Default(), store)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res
}
