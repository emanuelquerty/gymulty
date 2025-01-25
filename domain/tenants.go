package domain

import "time"

type Tenant struct {
	ID         int
	Name       string
	Email      string
	Password   string
	Subdomain  string
	Created_at time.Time
	Updated_at time.Time
}
