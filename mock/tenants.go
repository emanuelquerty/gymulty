package mock

import "github.com/emanuelquerty/multency/domain"

var _ domain.TenantStore = (*TenantStore)(nil)

type TenantStore struct {
	CreateTenantFn func(data domain.Tenant) (domain.Tenant, error)
}

func (t *TenantStore) CreateTenant(data domain.Tenant) (domain.Tenant, error) {
	return t.CreateTenantFn(data)
}
