package mcmodel

import (
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	ID                   int       `json:"id"`
	UUID                 string    `json:"uuid"`
	UsesUUID             string    `json:"uses_uuid"`
	UsesID               int       `json:"uses_id"`
	ProjectID            int       `json:"project_id"`
	Name                 string    `json:"name"`
	OwnerID              int       `json:"owner_id"`
	Path                 string    `json:"path"`
	DirectoryID          int       `json:"directory_id"`
	Size                 uint64    `json:"size"`
	Checksum             string    `json:"checksum"`
	MimeType             string    `json:"mime_type"`
	MediaTypeDescription string    `json:"media_type_description"`
	Current              bool      `json:"current"`
	Directory            *File     `json:"directory" gorm:"foreignKey:DirectoryID;references:ID"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
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

func (f File) ToUnderlyingFilePath(mcdir string) string {
	return filepath.Join(f.ToUnderlyingDirPath(mcdir), f.UUIDForPath())
}

func (f File) ToUnderlyingDirPath(mcdir string) string {
	uuidParts := strings.Split(f.UUIDForPath(), "-")
	return filepath.Join(mcdir, uuidParts[1][0:2], uuidParts[1][2:4])
}

func (f File) UUIDForPath() string {
	if f.UsesUUID != "" {
		return f.UsesUUID
	}

	return f.UUID
}
