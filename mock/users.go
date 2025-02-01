package mock

import (
	"github.com/emanuelquerty/multency/domain"
)

var _ domain.UserStore = (*UserStore)(nil)

type UserStore struct {
	GetUserByIDFn func(tenantID int, userID int) (domain.User, error)
	CreateUserFn  func(user domain.User) (domain.User, error)
	UpdateUserFn  func(id int, update domain.UserUpdate) (domain.User, error)
	DeleteByIDFn  func(id int) error
}

func (u *UserStore) GetUserByID(tenantID int, userID int) (domain.User, error) {
	return u.GetUserByIDFn(tenantID, userID)
}

func (u *UserStore) CreateUser(user domain.User) (domain.User, error) {
	return u.CreateUserFn(user)
}

func (u *UserStore) UpdateUser(id int, update domain.UserUpdate) (domain.User, error) {
	return u.UpdateUserFn(id, update)
}

func (u *UserStore) DeleteUserByID(id int) error {
	return u.DeleteByIDFn(id)
}
