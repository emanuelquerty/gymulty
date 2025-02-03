package http

import "github.com/emanuelquerty/gymulty/domain"

type Response[T any] struct {
	Success bool   `json:"success,omitempty"  bson:"success"`
	Count   int    `json:"count,omitempty"  bson:"count"`
	Type    string `json:"type,omitempty"  bson:"type"`
	Data    T      `json:"data,omitempty"  bson:"data"`
}

type TenantSignupResponse struct {
	Message string            `json:"message,omitempty"  bson:"message"`
	Tenant  domain.Tenant     `json:"tenant,omitempty"  bson:"tenant"`
	Admin   domain.PublicUser `json:"admin,omitempty"  bson:"admin"`
}
