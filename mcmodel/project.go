package mcmodel

import "time"

type Project struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	TeamID    int       `json:"team_id"`
	OwnerID   int       `json:"owner_id"`
	Owner     *User     `gorm:"foreignKey:OwnerID;references:ID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
