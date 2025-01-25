package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
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

func (u *User) HashPassword() error {
	password := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
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

// UserUpdates enables user to update one or more fields
// fields not nil are updated
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
