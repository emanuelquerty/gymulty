package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/emanuelquerty/gymulty/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	user := domain.User{
		ID:        1,
		TenantID:  1,
		FirstName: "Peter",
		LastName:  "Petrelli",
		Role:      "admin",
	}

	t.Run("returns 200 status code for existing user id ", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(ctx context.Context, tenantID int, userID int) (domain.User, error) {
			return user, nil
		}
		req := httptest.NewRequest("GET", "/api/tenants/1/users/1", nil)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 200
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns user with id 7", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(ctx context.Context, tenantID int, userID int) (domain.User, error) {
			found := user
			found.ID = 7
			return found, nil
		}
		req := httptest.NewRequest("GET", "/api/tenants/1/users/7", nil)
		res := newUserRequest(userStore, req)

		var got Response[[]domain.PublicUser]
		json.NewDecoder(res.Body).Decode(&got)

		want := Response[[]domain.PublicUser]{
			Count: 1,
			Data: []domain.PublicUser{
				{
					ID:        7,
					TenantID:  1,
					FirstName: "Peter",
					LastName:  "Petrelli",
					Role:      "admin",
				},
			},
		}
		assert.Equal(t, want, got, "users should be equal")
	})
	t.Run("returns user with id 2", func(t *testing.T) {
		user := domain.User{
			ID:        2,
			TenantID:  1,
			FirstName: "Bruce",
			LastName:  "Banner",
			Role:      "trainer",
		}
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(ctx context.Context, tenantID int, userID int) (domain.User, error) {
			return user, nil
		}
		req := httptest.NewRequest("GET", "/api/tenants/1/users/2", nil)
		res := newUserRequest(userStore, req)

		var got Response[[]domain.PublicUser]
		json.NewDecoder(res.Body).Decode(&got)

		want := Response[[]domain.PublicUser]{
			Count: 1,
			Data: []domain.PublicUser{
				{
					ID:        2,
					TenantID:  1,
					FirstName: "Bruce",
					LastName:  "Banner",
					Role:      "trainer",
				},
			},
		}
		assert.Equal(t, want, got, "users should be equal")
	})

	t.Run("returns 404 status code for non-existing user id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(ctx context.Context, tenantID int, userID int) (domain.User, error) {
			return domain.User{}, sql.ErrNoRows
		}

		req := httptest.NewRequest("GET", "/api/tenants/1/users/3", nil)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 404
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 400 status code for invalid user id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(ctx context.Context, tenantID int, userID int) (domain.User, error) {
			return domain.User{}, nil
		}

		req := httptest.NewRequest("GET", "/api/tenants/1/users/notValidID3", nil)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 400
		assert.Equal(t, want, got, "status codes should be equal")
	})

}

func TestCreateUser(t *testing.T) {
	user := domain.User{
		ID:        1,
		FirstName: "Leny",
		LastName:  "Jenkins",
		Email:     "ljenkins@email.com",
		Password:  "ReallyStrong21734bs",
		Role:      "admin",
	}

	t.Run("returns 201 status code", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(ctx context.Context, tenantID int, user domain.User) (domain.User, error) {
			return domain.User{}, nil
		}
		body, _ := json.Marshal(user)
		bodyBuff := bytes.NewBuffer(body)

		req := httptest.NewRequest(http.MethodPost, "/api/tenants/1/users", bodyBuff)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 201
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns newly created user", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(ctx context.Context, tenantID int, user domain.User) (domain.User, error) {
			return user, nil
		}

		body, _ := json.Marshal(user)
		bodyBuff := bytes.NewBuffer(body)

		req := httptest.NewRequest(http.MethodPost, "/api/tenants/1/users", bodyBuff)
		res := newUserRequest(userStore, req)

		var got Response[[]domain.PublicUser]
		json.NewDecoder(res.Body).Decode(&got)

		want := Response[[]domain.PublicUser]{
			Count: 1,
			Data: []domain.PublicUser{
				{
					ID:        user.ID,
					TenantID:  user.TenantID,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Role:      user.Role,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				},
			},
		}
		assert.Equal(t, want, got, "users should be equal")
	})

	t.Run("returns location header with full resource uri", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(ctx context.Context, tenantID int, user domain.User) (domain.User, error) {
			return user, nil
		}

		body, _ := json.Marshal(user)
		bodyBuff := bytes.NewBuffer(body)

		req := httptest.NewRequest(http.MethodPost, "/api/tenants/1/users", bodyBuff)
		res := newUserRequest(userStore, req)

		got := res.Header().Get("Location")
		want := "://example.com/api/tenants/1/users/1" // newly created resource always has mocked ID = 1
		assert.Equal(t, want, got, "urls should be equal")
	})
}

func TestUpdateUser(t *testing.T) {
	user := domain.User{
		ID:        3,
		TenantID:  1,
		FirstName: "Johnny",
		LastName:  "Presley",
		Email:     "jpres@email.com",
		Password:  "Very12SecuryPassword3245",
		Role:      "member",
	}

	t.Run("returns newly updated user", func(t *testing.T) {
		newFirstName := "Johnnyyyyyyy"
		newRole := "coach"
		update := domain.UserUpdate{
			FirstName: &newFirstName,
			Role:      &newRole,
		}
		userStore := new(mock.UserStore)
		userStore.UpdateUserFn = func(ctx context.Context, tenantID int, userID int, update domain.UserUpdate) (domain.User, error) {
			updatedUser := user
			updatedUser.FirstName = *update.FirstName
			updatedUser.Role = *update.Role
			return updatedUser, nil
		}

		want := Response[[]domain.PublicUser]{
			Count: 1,
			Data: []domain.PublicUser{
				{
					ID:        user.ID,
					TenantID:  user.TenantID,
					FirstName: *update.FirstName, // updated field
					LastName:  user.LastName,
					Role:      *update.Role, // updated field
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				},
			},
		}

		body, _ := json.Marshal(update)
		bodyBuff := bytes.NewBuffer(body)
		req := httptest.NewRequest(http.MethodPut, "/api/tenants/1/users/3", bodyBuff)
		res := newUserRequest(userStore, req)

		var got Response[[]domain.PublicUser]
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("returns 400 status code for invalid id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.UpdateUserFn = func(ctx context.Context, tenantID int, userID int, update domain.UserUpdate) (domain.User, error) {
			return domain.User{}, nil
		}
		req := httptest.NewRequest("PUT", "/api/tenants/1/users/notValidID8", nil)
		res := newUserRequest(userStore, req)

		got := res.Code
		want := 400
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 404 status code for non-existing user id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.UpdateUserFn = func(ctx context.Context, tenantID int, userID int, update domain.UserUpdate) (domain.User, error) {
			return domain.User{}, sql.ErrNoRows
		}
		req := httptest.NewRequest("PUT", "/api/tenants/1/users/27", nil)
		res := newUserRequest(userStore, req)

		got := res.Code
		want := 404
		assert.Equal(t, want, got, "status codes should be equal")
	})

}

func TestDelete(t *testing.T) {
	t.Run("returns 204 on success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/tenants/1/users/2", nil)
		userstore := new(mock.UserStore)
		userstore.DeleteByIDFn = func(ctx context.Context, tenantID int, userID int) error {
			return nil
		}

		res := newUserRequest(userstore, req)
		got := res.Code
		want := 204
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 404 for non-existing id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/tenants/1/users/5", nil)
		userstore := new(mock.UserStore)
		userstore.DeleteByIDFn = func(ctx context.Context, tenantID int, userID int) error {
			return sql.ErrNoRows
		}

		res := newUserRequest(userstore, req)
		got := res.Code
		want := 404
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 400 status code for invalid id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.DeleteByIDFn = func(ctx context.Context, tenantID int, userID int) error {
			return nil // this func is never called for this test, so return val here is irrelevant
		}
		req := httptest.NewRequest("DELETE", "/api/tenants/1/users/InvalidID", nil)
		res := newUserRequest(userStore, req)

		got := res.Code
		want := 400
		assert.Equal(t, want, got, "status codes should be equal")
	})
}

func TestGetAllUsers(t *testing.T) {
	users := []domain.User{
		{
			ID:        1,
			FirstName: "Leny",
			LastName:  "Jenkins",
			Email:     "ljenkins@email.com",
			Password:  "ReallyStrong21734bs",
			Role:      "admin",
		},
	}

	t.Run("returns all users on success given tenantID", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetAllUsersFn = func(ctx context.Context, tenantID int) ([]domain.User, error) {
			return users, nil
		}

		req := httptest.NewRequest("GET", "/api/tenants/1/users", nil)
		res := newUserRequest(userStore, req)

		var got Response[[]domain.PublicUser]
		json.NewDecoder(res.Body).Decode(&got)

		want := Response[[]domain.PublicUser]{
			Count: 1,
			Data:  MapToPublicUsers(users),
		}

		assert.Equal(t, want, got, "responses should match")
	})
}

func newUserRequest(userStore *mock.UserStore, req *http.Request) *httptest.ResponseRecorder {
	userHandler := NewUserHandler(slog.Default(), userStore)
	res := httptest.NewRecorder()
	userHandler.ServeHTTP(res, req)
	return res
}
