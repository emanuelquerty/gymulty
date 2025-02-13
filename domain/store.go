package domain

type Store interface {
	TenantStore
	UserStore
	ClassStore
}
