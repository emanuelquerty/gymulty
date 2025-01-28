package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emanuelquerty/multency/domain"
	"github.com/emanuelquerty/multency/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	user := domain.User{
		ID:        1,
		TenantID:  1,
		Firstname: "Peter",
		Lastname:  "Petrelli",
		Role:      "admin",
	}

	t.Run("returns 200 status code for existing user id ", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(id int) (domain.User, error) {
			return user, nil
		}
		req := httptest.NewRequest("GET", "/api/users/1", nil)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 200
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns user with id 7", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(id int) (domain.User, error) {
			found := user
			found.ID = 7
			return found, nil
		}
		req := httptest.NewRequest("GET", "/api/users/7", nil)
		res := newUserRequest(userStore, req)

		var got domain.PublicUser
		json.NewDecoder(res.Body).Decode(&got)

		want := domain.PublicUser{
			ID:        7,
			TenantID:  1,
			Firstname: "Peter",
			Lastname:  "Petrelli",
			Role:      "admin",
		}
		assert.Equal(t, want, got, "users should be equal")
	})
	t.Run("returns user with id 2", func(t *testing.T) {
		user := domain.User{
			ID:        2,
			TenantID:  1,
			Firstname: "Bruce",
			Lastname:  "Banner",
			Role:      "trainer",
		}
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(id int) (domain.User, error) {
			return user, nil
		}
		req := httptest.NewRequest("GET", "/api/users/2", nil)
		res := newUserRequest(userStore, req)

		var got domain.PublicUser
		json.NewDecoder(res.Body).Decode(&got)

		want := domain.PublicUser{
			ID:        2,
			TenantID:  1,
			Firstname: "Bruce",
			Lastname:  "Banner",
			Role:      "trainer",
		}
		assert.Equal(t, want, got, "users should be equal")
	})

	t.Run("returns 404 status code for non-existing user id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(id int) (domain.User, error) {
			return domain.User{}, errors.New("not found")
		}

		req := httptest.NewRequest("GET", "/api/users/3", nil)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 404
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 400 status code for invalid user id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.GetUserByIDFn = func(id int) (domain.User, error) {
			return domain.User{}, nil
		}

		req := httptest.NewRequest("GET", "/api/users/notValidID3", nil)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 400
		assert.Equal(t, want, got, "status codes should be equal")
	})

}

func TestCreateUser(t *testing.T) {
	user := domain.User{
		ID:        1,
		TenantID:  1,
		Firstname: "Leny",
		Lastname:  "Jenkins",
		Email:     "ljenkins@email.com",
		Password:  "ReallyStrong21734bs",
		Role:      "admin",
	}

	t.Run("returns 201 status code", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(user domain.User) (domain.User, error) {
			return domain.User{}, nil
		}
		body, _ := json.Marshal(user)
		bodyBuff := bytes.NewBuffer(body)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bodyBuff)
		res := newUserRequest(userStore, req)

		got, want := res.Code, 201
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns newly created user", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(user domain.User) (domain.User, error) {
			return user, nil
		}

		body, _ := json.Marshal(user)
		bodyBuff := bytes.NewBuffer(body)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bodyBuff)
		res := newUserRequest(userStore, req)

		var got domain.PublicUser
		json.NewDecoder(res.Body).Decode(&got)

		want := MapToPublicUser(user)
		assert.Equal(t, want, got, "users should be equal")
	})

	t.Run("returns location header with full resource uri", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.CreateUserFn = func(user domain.User) (domain.User, error) {
			return user, nil
		}

		body, _ := json.Marshal(user)
		bodyBuff := bytes.NewBuffer(body)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bodyBuff)
		res := newUserRequest(userStore, req)

		got := res.Header().Get("Location")
		want := "://example.com/users/1" // newly created resource always has mocked ID = 1
		assert.Equal(t, want, got, "urls should be equal")
	})
}

func TestUpdateUser(t *testing.T) {
	user := domain.User{
		ID:        3,
		TenantID:  1,
		Firstname: "Johnny",
		Lastname:  "Presley",
		Email:     "jpres@email.com",
		Password:  "Very12SecuryPassword3245",
		Role:      "member",
	}

	t.Run("returns newly updated user", func(t *testing.T) {
		newEmail := "Johnnypresley@email.com"
		newRole := "coach"
		update := domain.UserUpdate{
			Email: &newEmail,
			Role:  &newRole,
		}
		userStore := new(mock.UserStore)
		userStore.UpdateUserFn = func(id int, update domain.UserUpdate) (domain.User, error) {
			updatedUser := user
			updatedUser.Email = *update.Email
			updatedUser.Role = *update.Role
			return updatedUser, nil
		}
		want := user
		want.Email = *update.Email
		want.Role = *update.Role

		body, _ := json.Marshal(update)
		bodyBuff := bytes.NewBuffer(body)
		req := httptest.NewRequest(http.MethodPut, "/api/users/3", bodyBuff)
		res := newUserRequest(userStore, req)

		var got domain.User
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("returns 400 status code for invalid id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.UpdateUserFn = func(id int, update domain.UserUpdate) (domain.User, error) {
			return domain.User{}, nil
		}
		req := httptest.NewRequest("PUT", "/api/users/notValidID8", nil)
		res := newUserRequest(userStore, req)

		got := res.Code
		want := 400
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 404 status code for non-existing user id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.UpdateUserFn = func(id int, update domain.UserUpdate) (domain.User, error) {
			return domain.User{}, errors.New("not found")
		}
		req := httptest.NewRequest("PUT", "/api/users/27", nil)
		res := newUserRequest(userStore, req)

		got := res.Code
		want := 404
		assert.Equal(t, want, got, "status codes should be equal")
	})

}

func TestDelete(t *testing.T) {
	t.Run("returns 204 on success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/users/2", nil)
		userstore := new(mock.UserStore)
		userstore.DeleteByIDFn = func(id int) error {
			return nil
		}

		res := newUserRequest(userstore, req)
		got := res.Code
		want := 204
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 404 for non-existing id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/users/5", nil)
		userstore := new(mock.UserStore)
		userstore.DeleteByIDFn = func(id int) error {
			return errors.New("not found")
		}

		res := newUserRequest(userstore, req)
		got := res.Code
		want := 404
		assert.Equal(t, want, got, "status codes should be equal")
	})

	t.Run("returns 400 status code for invalid id", func(t *testing.T) {
		userStore := new(mock.UserStore)
		userStore.DeleteByIDFn = func(id int) error {
			return nil // this func is never called for this test, so return val here is irrelevant
		}
		req := httptest.NewRequest("DELETE", "/api/users/InvalidID", nil)
		res := newUserRequest(userStore, req)

		got := res.Code
		want := 400
		assert.Equal(t, want, got, "status codes should be equal")
	})
}

func newUserRequest(userStore *mock.UserStore, req *http.Request) *httptest.ResponseRecorder {
	userHandler := NewUserHandler(userStore)
	res := httptest.NewRecorder()
	userHandler.ServeHTTP(res, req)
	return res
}
