package mock

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
)

var _ domain.ClassStore = (*ClassStore)(nil)

type ClassStore struct {
	CreateClassFn func(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error)
}

func (c *ClassStore) CreateClass(ctx context.Context, tenantID int, class domain.Class) (domain.Class, error) {
	return c.CreateClassFn(ctx, tenantID, class)
}
