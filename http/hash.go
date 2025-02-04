package http

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
