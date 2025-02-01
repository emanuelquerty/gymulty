package mock

import (
	"github.com/emanuelquerty/multency/domain"
)

var _ domain.UserStore = (*UserStore)(nil)

type UserStore struct {
	GetUserByIDFn func(tenantID int, userID int) (domain.User, error)
	CreateUserFn  func(tenantID int, user domain.User) (domain.User, error)
	UpdateUserFn  func(tenantID int, userID int, update domain.UserUpdate) (domain.User, error)
	DeleteByIDFn  func(tenantID int, userID int) error
	GetAllUsersFn func(tenantID int) ([]domain.User, error)
}

func (u *UserStore) GetUserByID(tenantID int, userID int) (domain.User, error) {
	return u.GetUserByIDFn(tenantID, userID)
}

func (u *UserStore) CreateUser(tenantID int, user domain.User) (domain.User, error) {
	return u.CreateUserFn(tenantID, user)
}

func (u *UserStore) UpdateUser(tenantID int, userID int, update domain.UserUpdate) (domain.User, error) {
	return u.UpdateUserFn(tenantID, userID, update)
}

func (u *UserStore) DeleteUserByID(tenantID int, userID int) error {
	return u.DeleteByIDFn(tenantID, userID)
}

func (u *UserStore) GetAllUsers(tenantID int) ([]domain.User, error) {
	return u.GetAllUsersFn(tenantID)
}
