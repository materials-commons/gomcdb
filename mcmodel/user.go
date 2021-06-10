package mcmodel

import "time"

type User struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GlobusUser string `json:"globus_user"`
	ApiToken   string `json:"-"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
