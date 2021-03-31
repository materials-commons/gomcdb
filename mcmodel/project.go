package mcmodel

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	TeamID    int       `json:"team_id"`
	OwnerID   int       `json:"owner_id"`
	Owner     *User     `gorm:"foreignKey:OwnerID;references:ID"`
	FileTypes string    `json:"file_types"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p Project) GetFileTypes() (map[string]int, error) {
	var fileTypes map[string]int
	err := json.Unmarshal([]byte(p.FileTypes), &fileTypes)
	return fileTypes, err
}

func (p Project) ToFileTypeAsString(fileTypes map[string]int) (string, error) {
	b, err := json.Marshal(fileTypes)
	return string(b), err
}
