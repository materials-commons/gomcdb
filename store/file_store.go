package store

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/gomcdb/mcmodel"
	"gorm.io/gorm"
)

type FileStore struct {
	db       *gorm.DB
	mcfsRoot string
}

func NewFileStore(db *gorm.DB, mcfsRoot string) *FileStore {
	return &FileStore{db: db, mcfsRoot: mcfsRoot}
}

// UpdateMetadataForFileAndProject updates the metadata and project meta data for a file
func (s *FileStore) UpdateMetadataForFileAndProject(file *mcmodel.File, checksum string, projectID int, totalBytes int64) error {
	finfo, err := os.Stat(file.ToUnderlyingFilePath(s.mcfsRoot))
	if err != nil {
		log.Errorf("MarkFileReleased Stat %s failed: %s", file.ToUnderlyingFilePath(s.mcfsRoot), err)
		return err
	}

	return s.withTxRetry(func(tx *gorm.DB) error {
		// To set file as the current (ie viewable) version we first need to set all its previous
		// versions to have current set to false.
		err := tx.Model(&mcmodel.File{}).
			Where("directory_id = ?", file.DirectoryID).
			Where("name = ?", file.Name).
			Update("current", false).Error

		if err != nil {
			return err
		}

		// Now we can update the meta data on the current file. This includes, the size, current, and if there is
		// a new computed checksum, also update the checksum field.
		fileMetadata := mcmodel.File{
			Size:     uint64(finfo.Size()),
			Current:  true,
			Checksum: checksum,
		}

		if err := tx.Model(file).Updates(&fileMetadata).Error; err != nil {
			return err
		}

		var project mcmodel.Project

		if result := tx.Find(&project, projectID); result.Error != nil {
			return result.Error
		}

		return tx.Model(&project).Updates(&mcmodel.Project{Size: project.Size + totalBytes}).Error
	})
}

func (s *FileStore) CreateFile(name string, projectID, directoryID, ownerID int, mimeType string) (*mcmodel.File, error) {
	newFile := &mcmodel.File{
		ProjectID:   projectID,
		Name:        name,
		DirectoryID: directoryID,
		Size:        0,
		Checksum:    "",
		MimeType:    mimeType,
		OwnerID:     ownerID,
		Current:     false,
	}

	var err error

	if newFile.UUID, err = uuid.GenerateUUID(); err != nil {
		return nil, err
	}

	err = s.withTxRetry(func(tx *gorm.DB) error {
		return tx.Create(newFile).Error
	})

	return newFile, err
}

func (s *FileStore) FindDirByPath(projectID int, path string) (*mcmodel.File, error) {
	return findDirByPath(s.db, projectID, path)
}

func findDirByPath(db *gorm.DB, projectID int, path string) (*mcmodel.File, error) {
	var dir mcmodel.File
	err := db.Preload("Directory").
		Where("project_id = ?", projectID).
		Where("path = ?", path).
		Where("deleted_at IS NULL").
		Where("dataset_id IS NULL").
		First(&dir).Error
	if err != nil {
		//log.Errorf("Failed looking up directory in project %d, path %s: %s", projectID, path, err)
		return nil, err
	}

	return &dir, nil
}

func (s *FileStore) CreateDirectory(parentDirID, projectID, ownerID int, path, name string) (*mcmodel.File, error) {
	var dir mcmodel.File
	err := s.withTxRetry(func(tx *gorm.DB) error {
		err := tx.Where("path = ", path).
			Where("deleted_at IS NULL").
			Where("dataset_id IS NULL").
			Where("project_id = ?", projectID).
			Find(&dir).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			// directory already exists no need to create
			return nil
		}

		var project mcmodel.Project
		if result := tx.Find(&project, projectID); result.Error != nil {
			return result.Error
		}

		dir = mcmodel.File{
			OwnerID:              ownerID,
			MimeType:             "directory",
			MediaTypeDescription: "directory",
			DirectoryID:          parentDirID,
			Current:              true,
			Path:                 path,
			ProjectID:            projectID,
			Name:                 name,
		}

		if dir.UUID, err = uuid.GenerateUUID(); err != nil {
			return err
		}

		if err := tx.Create(&dir).Error; err != nil {
			return err
		}

		return tx.Model(&project).Updates(&mcmodel.Project{DirectoryCount: project.DirectoryCount + 1}).Error
	})

	return &dir, err
}

func (s *FileStore) CreateDirIfNotExists(parentDirID int, path, name string, projectID, ownerID int) (*mcmodel.File, error) {
	var (
		dir *mcmodel.File
		err error
	)

	err = s.withTxRetry(func(tx *gorm.DB) error {
		dir, err = findDirByPath(tx, projectID, path)
		if err == nil {
			// dir found
			return nil
		}

		dir = &mcmodel.File{
			OwnerID:              ownerID,
			MimeType:             "directory",
			MediaTypeDescription: "directory",
			DirectoryID:          parentDirID,
			Current:              true,
			Path:                 path,
			ProjectID:            projectID,
			Name:                 name,
		}

		if dir.UUID, err = uuid.GenerateUUID(); err != nil {
			return err
		}

		if err := tx.Create(dir).Error; err != nil {
			return err
		}

		project := mcmodel.Project{ID: projectID}

		return tx.Model(&project).Updates(&mcmodel.Project{DirectoryCount: project.DirectoryCount + 1}).Error
	})

	return dir, err
}

func (s *FileStore) ListDirectoryByPath(projectID int, path string) ([]mcmodel.File, error) {
	dir, err := s.FindDirByPath(projectID, path)
	if err != nil {
		return nil, err
	}

	var files []mcmodel.File

	err = s.db.Where("directory_id = ?", dir.ID).
		Where("project_id", projectID).
		Where("deleted_at IS NULL").
		Where("dataset_id IS NULL").
		Where("current = true").
		Find(&files).Error

	return files, err
}

// FindOrCreateDirPath will create all entries in the directory path if the path doesn't exist
func (s *FileStore) FindOrCreateDirPath(projectID, ownerID int, path string) (*mcmodel.File, error) {
	dir, err := s.FindDirByPath(projectID, path)
	if err == nil {
		// If we are here then the path was found, and we have nothing left to do.
		return dir, nil
	}

	// If we are here then directory wasn't found. At this point we don't know how many levels deep we have
	// to create directories. The common case is that this directory doesn't exist, but the parent does. Let's
	// check that case since it saves us a lot of work.
	parentPath := filepath.Dir(path)
	parentDir, err := s.FindDirByPath(projectID, parentPath)
	if err == nil {
		// Ok, the parent exists, so just create the child of the parent (ie, the complete path) and return
		// the created directory.
		return s.CreateDirectory(parentDir.ID, projectID, ownerID, path, filepath.Base(path))
	}

	// If we are here then the path didn't exist and the parent didn't exist so now we are going to traverse
	// upwards constructing as we go. The way we do this is to split the path, retrieve the root, and then just
	// start appending each entry of the path on, checking if it exists and if it doesn't then create it.

	// Start with root and then go from there
	parentDir, err = s.FindDirByPath(projectID, "/")
	if err != nil {
		return nil, err
	}

	pathParts := filepath.SplitList(path)
	currentPath := "/"
	for _, pathPart := range pathParts[1:] {
		currentPath = filepath.Join(currentPath, pathPart)
		dir, err = s.CreateDirIfNotExists(parentDir.ID, currentPath, filepath.Base(currentPath), projectID, ownerID)
		if err != nil {
			return nil, err
		}
		parentDir = dir
	}

	return dir, nil
}

func (s *FileStore) FindFileByPath(projectID int, path string) (*mcmodel.File, error) {
	dirPath := filepath.Dir(path)
	dir, err := s.FindDirByPath(projectID, dirPath)
	if err != nil {
		return nil, err
	}

	var file mcmodel.File
	err = s.db.Preload("Directory").
		Where("directory_id = ?", dir.ID).
		Where("name = ?", filepath.Base(path)).
		Where("deleted_at IS NULL").
		Where("dataset_id IS NULL").
		Where("current = ?", true).
		First(&file).Error
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *FileStore) UpdateFileUses(file *mcmodel.File, uuid string, fileID int) error {
	return s.withTxRetry(func(tx *gorm.DB) error {
		return tx.Model(file).Updates(mcmodel.File{
			UsesUUID: uuid,
			UsesID:   fileID,
		}).Error
	})
}

func (s *FileStore) PointAtExistingIfExists(file *mcmodel.File) (bool, error) {
	switched := false // Set to true in withTxRetry if an existing file with same checksum is found
	err := s.withTxRetry(func(tx *gorm.DB) error {
		var matched mcmodel.File
		err := tx.Where("checksum = ?", file.Checksum).
			Where("deleted_at IS NULL").
			Where("id <> ?", file.ID).
			First(&matched).Error
		if err == nil {
			// found a match
			switched = true
			usesUUID := matched.UUIDForUses()
			usesID := matched.IDForUses()
			return tx.Model(file).Updates(mcmodel.File{
				UsesUUID: usesUUID,
				UsesID:   usesID,
			}).Error
		}
		return nil
	})

	return switched, err
}

func (s *FileStore) withTxRetry(fn func(tx *gorm.DB) error) error {
	return WithTxRetryDefault(fn, s.db)
}
