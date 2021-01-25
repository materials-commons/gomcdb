package mcmodel

import (
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	ID          int       `json:"id"`
	UUID        string    `json:"string"`
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	OwnerID     int       `json:"owner_id"`
	Path        string    `json:"path"`
	DirectoryID int       `json:"directory_id"`
	Size        uint64    `json:"size"`
	Checksum    string    `json:"checksum"`
	MimeType    string    `json:"mime_type"`
	Current     bool      `json:"current"`
	Directory   *File     `json:"directory" gorm:"foreignKey:DirectoryID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (File) TableName() string {
	return "files"
}

func (f File) IsFile() bool {
	return f.MimeType != "directory"
}

func (f File) IsDir() bool {
	return f.MimeType == "directory"
}

func (f File) FullPath() string {
	if f.IsDir() {
		return f.Path
	}

	// f is a file and not a directory
	if f.Directory.Path == "/" {
		return f.Directory.Path + f.Name
	}

	return f.Directory.Path + "/" + f.Name
}

func (f File) ToPath(mcdir string) string {
	return filepath.Join(f.ToDirPath(mcdir), f.UUID)
}

func (f File) ToDirPath(mcdir string) string {
	uuidParts := strings.Split(f.UUID, "-")
	return filepath.Join(mcdir, uuidParts[1][0:2], uuidParts[1][2:4])
}
