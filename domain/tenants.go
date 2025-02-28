package domain

import (
	"context"
	"time"
)

type Tenant struct {
	ID           int
	BusinessName string
	Subdomain    string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TenantRequestBody struct {
	BusinessName string `json:"business_name,omitempty"  bson:"business_name"`
	Subdomain    string `json:"subdomain,omitempty"  bson:"subdomain"`
	FirstName    string `json:"first_name,omitempty"  bson:"first_name"`
	LastName     string `json:"last_name,omitempty"  bson:"last_name"`
	Email        string `json:"email,omitempty"  bson:"email"`
	Password     string `json:"password,omitempty"  bson:"password"`
}

type TenantStore interface {
	CreateTenant(ctx context.Context, data Tenant) (Tenant, error)
}
