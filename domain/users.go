package domain

import "time"

type User struct {
	ID           int       `json:"id,omitempty"  bson:"id"`
	TenantID     int       `json:"tenant_id,omitempty"  bson:"tenant_id"`
	Firstname    string    `json:"firstname,omitempty"  bson:"firstname"`
	Lastname     string    `json:"lastname,omitempty"  bson:"lastname"`
	Email        string    `json:"email,omitempty"  bson:"email"`
	PasswordHash string    `json:"password_hash,omitempty"  bson:"password_hash"`
	Role         string    `json:"role,omitempty"  bson:"role"`
	CreatedAt    time.Time `json:"created_at,omitempty"  bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"  bson:"updated_at"`
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

type UserStore interface {
	GetUserByID(id int) (User, error)
}
