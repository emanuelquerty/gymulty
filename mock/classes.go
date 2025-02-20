package mock

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
)

var _ domain.ClassStore = (*ClassStore)(nil)

type ClassStore struct {
	CreateClassFn     func(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error)
	GetClassByIDFn    func(ctx context.Context, tenantID int, classID int) (domain.Class, error)
	DeleteClassByIDFn func(ctx context.Context, tenantID int, classID int) error
	GetAllClassesFn   func(ctx context.Context, tenantID int) ([]domain.Class, error)
}

func (c *ClassStore) CreateClass(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error) {
	return c.CreateClassFn(ctx, tenantID, class)
}

func (c *ClassStore) GetClassByID(ctx context.Context, tenantID int, classID int) (domain.Class, error) {
	return c.GetClassByIDFn(ctx, tenantID, classID)
}

func (c *ClassStore) DeleteClassByID(ctx context.Context, tenantID int, classID int) error {
	return c.DeleteClassByIDFn(ctx, tenantID, classID)
}

func (c *ClassStore) GetAllClasses(ctx context.Context, tenantID int) ([]domain.Class, error) {
	return c.GetAllClassesFn(ctx, tenantID)
}
