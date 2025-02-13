package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/emanuelquerty/gymulty/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateClass(t *testing.T) {
	class := domain.Class{
		ID:          1,
		TenantID:    6,
		TrainerID:   5,
		Name:        "Yoga Session",
		Description: "An amazing class that will bring you to a complete state of relaxation",
		Capacity:    18,
		StartsAt:    time.Now().AddDate(0, 0, 18),
		EndsAt:      time.Now().AddDate(0, 0, 18).Add(1 * time.Hour),
	}

	store := new(mock.ClassStore)
	store.CreateClassFn = func(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error) {
		return class, nil
	}

	t.Run("returns location header with resource uri of newly created class", func(t *testing.T) {
		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("GET", "/api/tenants/6/classes", buf)

		res := NewClassRequest(req, store)
		want := fmt.Sprintf("%s://%s/api/tenants/6/classes/%d", req.URL.Scheme, req.Host, class.ID)

		got := res.Header().Get("Location")
		assert.Equal(t, want, got, "uri in location header should match")
	})

	t.Run("returns newly created class", func(t *testing.T) {
		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("GET", "/api/tenants/6/classes", buf)

		res := NewClassRequest(req, store)
		want := class

		var got domain.Class
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "classes should match")
	})

	t.Run("returns 400 status code on invalid tenant id", func(t *testing.T) {

		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("GET", "/api/tenants/invalid324/classes", buf)

		store := new(mock.ClassStore)
		store.CreateClassFn = func(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error) {
			return domain.Class{}, nil
		}

		res := NewClassRequest(req, store)
		want := 400
		got := res.Code
		assert.Equal(t, want, got, "classes should match")
	})

	t.Run("returns 404 status code on non-existing tenant id", func(t *testing.T) {

		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("GET", "/api/tenants/99999/classes", buf)

		store := new(mock.ClassStore)
		store.CreateClassFn = func(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error) {
			return domain.Class{}, sql.ErrNoRows
		}

		res := NewClassRequest(req, store)
		want := 404
		got := res.Code
		assert.Equal(t, want, got, "classes should match")
	})
}

func NewClassRequest(req *http.Request, store domain.ClassStore) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	handler := NewClassHandler(slog.Default(), store)
	handler.ServeHTTP(res, req)
	return res
}
