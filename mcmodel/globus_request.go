package mcmodel

import "time"

type GlobusRequest struct {
	ID               int       `json:"id"`
	UUID             string    `json:"string"`
	ProjectID        int       `json:"project_id"`
	Name             string    `json:"name"`
	Pid              int       `json:"pid"`
	State            string    `json:"state"`
	OwnerID          int       `json:"owner_id"`
	Owner            *User     `gorm:"foreignKey:OwnerID;references:ID"`
	Path             string    `json:"path"`
	GlobusAclID      string    `json:"globus_acl_id"`
	GlobusPath       string    `json:"globus_path"`
	GlobusIdentityID string    `json:"globus_identity_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (GlobusRequest) TableName() string {
	return "globus_requests"
}
