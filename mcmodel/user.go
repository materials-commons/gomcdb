package mcmodel

import "time"

type User struct {
	ID         int
	UUID       string
	Name       string
	Email      string
	GlobusUser string
	ApiToken   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
