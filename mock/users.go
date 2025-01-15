package mock

import (
	"errors"

	"github.com/emanuelquerty/multency/domain"
)

// user1 := domain.PublicUser{
// 	ID:        1,
// 	TenantID:  1,
// 	Firstname: "Peter",
// 	Lastname:  "Petrelli",
// 	Role:      "admin",
// }

// user2 := domain.PublicUser{
// 	ID:        2,
// 	TenantID:  1,
// 	Firstname: "Bruce",
// 	Lastname:  "Benner",
// 	Role:      "trainer",
// }

type UserStore struct {
	users map[int]domain.User
}

func NewUserStore(users map[int]domain.User) *UserStore {
	return &UserStore{
		users: users,
	}
}
func (u *UserStore) GetUserByID(id int) (domain.User, error) {
	user, ok := u.users[id]
	if !ok {
		return domain.User{}, errors.New("")
	}
	return user, nil
}
