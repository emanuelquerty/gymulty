package domain

import (
	"time"
)

type User struct {
	ID        int       `json:"id,omitempty"  bson:"id"`
	TenantID  int       `json:"tenant_id,omitempty"  bson:"tenant_id"`
	Firstname string    `json:"firstname,omitempty"  bson:"firstname"`
	Lastname  string    `json:"lastname,omitempty"  bson:"lastname"`
	Email     string    `json:"email,omitempty"  bson:"email"`
	Password  string    `json:"password,omitempty"  bson:"password"`
	Role      string    `json:"role,omitempty"  bson:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"  bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"  bson:"updated_at"`
}

type PublicUser struct {
	ID        int       `json:"id,omitempty"  bson:"id"`
	TenantID  int       `json:"tenant_id,omitempty"  bson:"tenant_id"`
	Firstname string    `json:"firstname,omitempty"  bson:"firstname"`
	Lastname  string    `json:"lastname,omitempty"  bson:"lastname"`
	Role      string    `json:"role,omitempty"  bson:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"  bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"  bson:"updated_at"`
}

// UserUpdates contains all fields the user can manually update
// All fields are nil to make it easy to mass update
// Fields not nil are updated through reflection
// this makes it easy to update 1 or more fields at once
type UserUpdate struct {
	Firstname *string `json:"firstname,omitempty"  bson:"firstname"`
	Lastname  *string `json:"lastname,omitempty"  bson:"lastname"`
	Email     *string `json:"email,omitempty"  bson:"email"`
	Role      *string `json:"role,omitempty"  bson:"role"`
	Password  *string `json:"password,omitempty"  bson:"password"`
}

type UserStore interface {
	GetUserByID(id int) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(id int, updates UserUpdate) (User, error)
}
