package mock

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
)

var _ domain.TenantStore = (*TenantStore)(nil)

type TenantStore struct {
	CreateTenantFn func(ctx context.Context, data domain.Tenant) (domain.Tenant, error)
}

func (t *TenantStore) CreateTenant(ctx context.Context, data domain.Tenant) (domain.Tenant, error) {
	return t.CreateTenantFn(ctx, data)
}
