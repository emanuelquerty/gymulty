package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id,omitempty"  bson:"id"`
	TenantID  int       `json:"tenant_id,omitempty"  bson:"tenant_id"`
	FirstName string    `json:"first_name,omitempty"  bson:"firstname"`
	LastName  string    `json:"last_name,omitempty"  bson:"lastname"`
	Email     string    `json:"email,omitempty"  bson:"email"`
	Password  string    `json:"password,omitempty"  bson:"password"`
	Role      string    `json:"role,omitempty"  bson:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"  bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"  bson:"updated_at"`
}

type PublicUser struct {
	ID        int       `json:"id,omitempty"  bson:"id"`
	TenantID  int       `json:"tenant_id,omitempty"  bson:"tenant_id"`
	FirstName string    `json:"first_name,omitempty"  bson:"firstname"`
	LastName  string    `json:"last_name,omitempty"  bson:"lastname"`
	Role      string    `json:"role,omitempty"  bson:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"  bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"  bson:"updated_at"`
}

// UserUpdates enables user to update one or more fields
// fields not nil are updated
type UserUpdate struct {
	FirstName *string `json:"first_name,omitempty"  bson:"first_name"`
	LastName  *string `json:"last_name,omitempty"  bson:"last_name"`
	Email     *string `json:"email,omitempty"  bson:"email"`
	Role      *string `json:"role,omitempty"  bson:"role"`
	Password  *string `json:"password,omitempty"  bson:"password"`
}

type UserStore interface {
	CreateUser(ctx context.Context, tenantID int, user User) (User, error)
	GetUserByID(ctx context.Context, tenantID int, userID int) (User, error)
	UpdateUser(ctx context.Context, tenantID int, userID int, updates UserUpdate) (User, error)
	DeleteUserByID(ctx context.Context, tenantID int, userID int) error
	GetAllUsers(ctx context.Context, tenantID int) ([]User, error)
}
