package http

import "github.com/emanuelquerty/multency/domain"

type TenantSignupResponse struct {
	Message string
	Tenant  domain.Tenant
	User    domain.PublicUser
}
