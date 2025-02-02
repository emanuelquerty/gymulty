package mock

import (
	"context"

	"github.com/emanuelquerty/multency/domain"
)

var _ domain.UserStore = (*UserStore)(nil)

type UserStore struct {
	GetUserByIDFn func(ctx context.Context, tenantID int, userID int) (domain.User, error)
	CreateUserFn  func(ctx context.Context, tenantID int, user domain.User) (domain.User, error)
	UpdateUserFn  func(ctx context.Context, tenantID int, userID int, update domain.UserUpdate) (domain.User, error)
	DeleteByIDFn  func(ctx context.Context, tenantID int, userID int) error
	GetAllUsersFn func(ctx context.Context, tenantID int) ([]domain.User, error)
}

func (u *UserStore) GetUserByID(ctx context.Context, tenantID int, userID int) (domain.User, error) {
	return u.GetUserByIDFn(ctx, tenantID, userID)
}

func (u *UserStore) CreateUser(ctx context.Context, tenantID int, user domain.User) (domain.User, error) {
	return u.CreateUserFn(ctx, tenantID, user)
}

func (u *UserStore) UpdateUser(ctx context.Context, tenantID int, userID int, update domain.UserUpdate) (domain.User, error) {
	return u.UpdateUserFn(ctx, tenantID, userID, update)
}

func (u *UserStore) DeleteUserByID(ctx context.Context, tenantID int, userID int) error {
	return u.DeleteByIDFn(ctx, tenantID, userID)
}

func (u *UserStore) GetAllUsers(ctx context.Context, tenantID int) ([]domain.User, error) {
	return u.GetAllUsersFn(ctx, tenantID)
}
