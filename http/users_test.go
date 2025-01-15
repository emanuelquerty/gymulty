package http

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/emanuelquerty/multency/domain"
	"github.com/emanuelquerty/multency/mock"
)

func TestGetUser(t *testing.T) {
	t.Run("return user with id 1", func(t *testing.T) {
		users := map[int]domain.User{
			1: {
				ID:        1,
				TenantID:  1,
				Firstname: "Peter",
				Lastname:  "Petrelli",
				Role:      "admin",
			},
		}
		res := newUserRequest(users, "GET", "/api/users/1", nil)

		var got domain.PublicUser
		json.NewDecoder(res.Body).Decode(&got)

		want := domain.PublicUser{
			ID:        1,
			TenantID:  1,
			Firstname: "Peter",
			Lastname:  "Petrelli",
			Role:      "admin",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("return user with id 2", func(t *testing.T) {
		users := map[int]domain.User{
			2: {
				ID:        2,
				TenantID:  1,
				Firstname: "Bruce",
				Lastname:  "Banner",
				Role:      "trainer",
			},
		}
		res := newUserRequest(users, "GET", "/api/users/2", nil)

		var got domain.PublicUser
		json.NewDecoder(res.Body).Decode(&got)

		want := domain.PublicUser{
			ID:        2,
			TenantID:  1,
			Firstname: "Bruce",
			Lastname:  "Banner",
			Role:      "trainer",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("return 404 not found", func(t *testing.T) {
		users := make(map[int]domain.User)
		res := newUserRequest(users, "GET", "/api/users/3", nil)

		if got, want := res.Code, 404; got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("return 400 bad request", func(t *testing.T) {
		users := make(map[int]domain.User)
		res := newUserRequest(users, "GET", "/api/users/notValidID3", nil)

		if got, want := res.Code, 400; got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

}

func newUserRequest(users map[int]domain.User, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	userStore := mock.NewUserStore(users)
	userHandler := NewUserHandler(userStore)

	req := httptest.NewRequest(method, url, body)
	res := httptest.NewRecorder()

	userHandler.ServeHTTP(res, req)
	return res
}
