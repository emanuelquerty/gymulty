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

	t.Run("creates a new user, returning location header with resource uri", func(t *testing.T) {
		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/api/tenants/6/classes", buf)

		res := NewClassRequest(req, store)
		want := fmt.Sprintf("%s://%s/api/tenants/6/classes/%d", req.URL.Scheme, req.Host, class.ID)

		got := res.Header().Get("Location")
		assert.Equal(t, want, got, "uri in location header should match")
	})

	t.Run("creates a new user, returning newly created user", func(t *testing.T) {
		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/api/tenants/6/classes", buf)

		res := NewClassRequest(req, store)
		want := Response[[]domain.Class]{Count: 1, Data: []domain.Class{class}}

		var got Response[[]domain.Class]
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "classes should match")
	})

	t.Run("creates a new user, returning 201 status code", func(t *testing.T) {
		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/api/tenants/6/classes", buf)

		res := NewClassRequest(req, store)
		want := 201
		got := res.Code

		assert.Equal(t, want, got, "status codes should match")
	})

	t.Run("returns 400 status code on invalid tenant id", func(t *testing.T) {

		body, _ := json.Marshal(class)
		buf := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/api/tenants/invalid324/classes", buf)

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
		req := httptest.NewRequest("POST", "/api/tenants/99999/classes", buf)

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

func TestGetClassByID(t *testing.T) {
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
	t.Run("returns class with id 1", func(t *testing.T) {
		store := new(mock.ClassStore)
		store.GetClassByIDFn = func(ctx context.Context, tenantID, classID int) (domain.Class, error) {
			return class, nil
		}

		req := httptest.NewRequest("GET", "/api/tenants/6/classes/1", nil)
		res := NewClassRequest(req, store)

		want := Response[[]domain.Class]{Count: 1, Data: []domain.Class{class}}
		var got Response[[]domain.Class]
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "classes should match")
	})

	t.Run("returns 400 on invalid tenant/class id", func(t *testing.T) {
		store := new(mock.ClassStore)
		store.GetClassByIDFn = func(ctx context.Context, tenantID, classID int) (domain.Class, error) {
			return domain.Class{}, nil
		}

		//Invalid tenant id
		want := 400

		req := httptest.NewRequest("GET", "/api/tenants/INvalid6283id/classes/1", nil)
		res := NewClassRequest(req, store)
		got := res.Code
		assert.Equal(t, want, got, "status code should match")

		//Invalid class id
		req = httptest.NewRequest("GET", "/api/tenants/6/classes/invalid_classID1", nil)
		res = NewClassRequest(req, store)
		got = res.Code
		assert.Equal(t, want, got, "status code should match")
	})

	t.Run("returns 404 on non-existing tenant/class id", func(t *testing.T) {
		store := new(mock.ClassStore)
		store.GetClassByIDFn = func(ctx context.Context, tenantID, classID int) (domain.Class, error) {
			return domain.Class{}, sql.ErrNoRows
		}

		// Non-exisitng tenant id
		want := 404

		req := httptest.NewRequest("GET", "/api/tenants/637/classes/1", nil)
		res := NewClassRequest(req, store)
		got := res.Code
		assert.Equal(t, want, got, "status code should match")

		// Non-existing class id
		req = httptest.NewRequest("GET", "/api/tenants/6/classes/981", nil)
		res = NewClassRequest(req, store)
		got = res.Code
		assert.Equal(t, want, got, "status code should match")
	})
}

func TestDeleteClassByID(t *testing.T) {
	t.Run("delete class with id 3, returning 204 on success", func(t *testing.T) {
		store := new(mock.ClassStore)
		store.DeleteClassByIDFn = func(ctx context.Context, tenantID, classID int) error {
			return nil
		}
		want := 204

		req := httptest.NewRequest("DELETE", "/api/tenants/6/classes/1", nil)
		res := NewClassRequest(req, store)
		got := res.Code
		assert.Equal(t, want, got, "status code should match")
	})

	t.Run("delete class with invalid class id, returning 400 status code", func(t *testing.T) {
		store := new(mock.ClassStore)
		store.DeleteClassByIDFn = func(ctx context.Context, tenantID, classID int) error {
			return nil
		}
		want := 400

		req := httptest.NewRequest("DELETE", "/api/tenants/6/classes/INValidID10293", nil)
		res := NewClassRequest(req, store)
		got := res.Code
		assert.Equal(t, want, got, "status code should match")
	})

	t.Run("delete class with invalid tenant id, returning 400 status code", func(t *testing.T) {
		store := new(mock.ClassStore)
		store.DeleteClassByIDFn = func(ctx context.Context, tenantID, classID int) error {
			return nil
		}
		want := 400

		req := httptest.NewRequest("DELETE", "/api/tenants/NotValidID123/classes/1", nil)
		res := NewClassRequest(req, store)
		got := res.Code
		assert.Equal(t, want, got, "status code should match")
	})
}

func NewClassRequest(req *http.Request, store domain.ClassStore) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	handler := NewClassHandler(slog.Default(), store)
	handler.ServeHTTP(res, req)
	return res
}
