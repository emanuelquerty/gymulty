package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/emanuelquerty/multency/domain"
	"github.com/emanuelquerty/multency/mock"
)

func TestGetUser(t *testing.T) {

	t.Run("returns 200 status code for existing user id ", func(t *testing.T) {
		users := map[int]domain.User{
			1: {
				ID:        1,
				TenantID:  1,
				Firstname: "Peter",
				Lastname:  "Petrelli",
				Role:      "admin",
			},
		}
		req := httptest.NewRequest("GET", "/api/users/1", nil)
		res := newUserRequest(users, req)

		if got, want := res.Code, 200; got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("returns user with id 1", func(t *testing.T) {
		users := map[int]domain.User{
			1: {
				ID:        1,
				TenantID:  1,
				Firstname: "Peter",
				Lastname:  "Petrelli",
				Role:      "admin",
			},
		}
		req := httptest.NewRequest("GET", "/api/users/1", nil)
		res := newUserRequest(users, req)

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
	t.Run("returns user with id 2", func(t *testing.T) {
		users := map[int]domain.User{
			2: {
				ID:        2,
				TenantID:  1,
				Firstname: "Bruce",
				Lastname:  "Banner",
				Role:      "trainer",
			},
		}
		req := httptest.NewRequest("GET", "/api/users/2", nil)
		res := newUserRequest(users, req)

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

	t.Run("returns 404 status code", func(t *testing.T) {
		users := make(map[int]domain.User)

		req := httptest.NewRequest("GET", "/api/users/3", nil)
		res := newUserRequest(users, req)

		if got, want := res.Code, 404; got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("returns 400 status code", func(t *testing.T) {
		users := make(map[int]domain.User)

		req := httptest.NewRequest("GET", "/api/users/notValidID3", nil)
		res := newUserRequest(users, req)

		if got, want := res.Code, 400; got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

}

func TestCreateUser(t *testing.T) {
	storeData := make(map[int]domain.User)
	user := domain.User{
		TenantID:  1,
		Firstname: "Leny",
		Lastname:  "Jenkins",
		Email:     "ljenkins@email.com",
		Password:  "ReallyStrong21734bs",
		Role:      "admin",
	}

	t.Run("returns 201 with newly created user", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(user)
		reqBody := bytes.NewBuffer(bodyBytes)

		req := httptest.NewRequest(http.MethodPost, "/api/users", reqBody)
		res := newUserRequest(storeData, req)

		if got, want := res.Code, 201; got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		var got domain.PublicUser
		json.NewDecoder(res.Body).Decode(&got)

		user.ID = got.ID // user is what we want the server to return on creation, so just update with id

		if want := MapToPublicUser(user); !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("returns location header with full resource uri", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(user)
		reqBody := bytes.NewBuffer(bodyBytes)

		req := httptest.NewRequest(http.MethodPost, "/api/users", reqBody)
		res := newUserRequest(storeData, req)

		got := res.Header().Get("Location")
		want := "/api/users/1" // newly created resource always has mocked ID = 1
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}

func newUserRequest(users map[int]domain.User, req *http.Request) *httptest.ResponseRecorder {
	userStore := mock.NewUserStore(users)
	userHandler := NewUserHandler(userStore)

	res := httptest.NewRecorder()
	userHandler.ServeHTTP(res, req)
	return res
}
