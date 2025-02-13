package domain

import (
	"context"
	"time"
)

type Class struct {
	ID          int       `json:"id,omitempty"  bson:"id"`
	TenantID    int       `json:"tenant_id,omitempty"  bson:"tenant_id"`
	TrainerID   int       `json:"trainer_id,omitempty"  bson:"trainer_id"`
	Name        string    `json:"name,omitempty"  bson:"name"`
	Description string    `json:"description,omitempty"  bson:"descripption"`
	Capacity    int       `json:"capacity,omitempty"  bson:"capacity"`
	StartsAt    time.Time `json:"starts_at,omitempty"  bson:"starts_at"`
	EndsAt      time.Time `json:"ends_at,omitempty"  bson:"ends_at"`
	CreatedAt   time.Time `json:"created_at,omitempty"  bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"  bson:"updated_at"`
}

type ClassStore interface {
	CreateClass(ctx context.Context, tenantID int, class Class) (Class, error)
}
