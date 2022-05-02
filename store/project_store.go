package store

import (
	"github.com/materials-commons/gomcdb/mcmodel"
	"gorm.io/gorm"
)

type ProjectStore struct {
	db *gorm.DB
}

func NewProjectStore(db *gorm.DB) *ProjectStore {
	return &ProjectStore{db: db}
}

func (s *ProjectStore) FindProject(projectID int) (*mcmodel.Project, error) {
	var project mcmodel.Project
	err := s.db.Find(&project, projectID).Error
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectStore) GetProjectBySlug(slug string) (*mcmodel.Project, error) {
	var project mcmodel.Project
	if err := s.db.Where("slug = ?", slug).First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectStore) GetProjectsForUser(userID int) (error, []mcmodel.Project) {
	var projects []mcmodel.Project

	err := s.db.Where("team_id in (select team_id from team2admin where user_id = ?)", userID).
		Or("team_id in (select team_id from team2member where user_id = ?)", userID).
		Find(&projects).Error
	return err, projects
}

func (s *ProjectStore) UpdateProjectSizeAndFileCount(projectID int, size int64, fileCount int) error {
	return s.withTxRetry(func(tx *gorm.DB) error {
		var p mcmodel.Project
		// Get latest project
		if result := tx.Find(&p, projectID); result.Error != nil {
			return result.Error
		}

		return tx.Model(&p).Updates(mcmodel.Project{
			FileCount: p.FileCount + fileCount,
			Size:      p.Size + size,
		}).Error
	})
}

func (s *ProjectStore) UpdateProjectDirectoryCount(projectID int, directoryCount int) error {
	return s.withTxRetry(func(tx *gorm.DB) error {
		var p mcmodel.Project
		// Get latest project
		if result := tx.Find(&p, projectID); result.Error != nil {
			return result.Error
		}

		return tx.Model(&p).Updates(mcmodel.Project{
			DirectoryCount: p.DirectoryCount + directoryCount,
		}).Error
	})
}

func (s *ProjectStore) UserCanAccessProject(userID, projectID int) bool {
	var project mcmodel.Project

	if err := s.db.Find(&project, projectID).Error; err != nil {
		return false
	}

	if project.OwnerID == userID {
		return true
	}

	var userCount int64
	s.db.Table("team2member").
		Where("user_id = ?", userID).
		Where("team_id = ?", project.TeamID).
		Count(&userCount)

	if userCount != 0 {
		return true
	}

	s.db.Table("team2admin").
		Where("user_id = ?", userID).
		Where("team_id = ?", project.TeamID).
		Count(&userCount)

	if userCount != 0 {
		return true
	}

	return false
}

func (s *ProjectStore) withTxRetry(fn func(tx *gorm.DB) error) error {
	return WithTxRetryDefault(fn, s.db)
}
