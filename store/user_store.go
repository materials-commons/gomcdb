package store

import (
	"github.com/materials-commons/gomcdb/mcmodel"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetUsersWithGlobusAccount() (error, []mcmodel.User) {
	var users []mcmodel.User
	result := s.db.Where("globus_user is not null").Find(&users)
	return result.Error, users
}

func (s *UserStore) GetUserBySlug(slug string) (*mcmodel.User, error) {
	var user mcmodel.User
	if err := s.db.Where("slug = ?", slug).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
