package mock

import (
	"errors"

	"github.com/emanuelquerty/multency/domain"
)

type UserStore struct {
	users map[int]domain.User
}

func NewUserStore(users map[int]domain.User) *UserStore {
	return &UserStore{
		users: users,
	}
}
func (u *UserStore) GetUserByID(id int) (domain.User, error) {
	user, ok := u.users[id]
	if !ok {
		return domain.User{}, errors.New("")
	}
	return user, nil
}

func (u *UserStore) CreateUser(user domain.User) (domain.User, error) {
	user.ID = 1
	u.users[1] = user
	return u.users[1], nil
}
