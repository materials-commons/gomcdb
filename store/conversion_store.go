package store

import (
	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/gomcdb/mcmodel"
	"gorm.io/gorm"
)

type ConversionStore struct {
	db *gorm.DB
}

func NewConversionStore(db *gorm.DB) *ConversionStore {
	return &ConversionStore{db: db}
}

func (s *ConversionStore) AddFileToConvert(file *mcmodel.File) (*mcmodel.Conversion, error) {
	var err error

	c := &mcmodel.Conversion{
		ProjectID: file.ProjectID,
		FileID:    file.ID,
		OwnerID:   file.OwnerID,
	}

	if c.UUID, err = uuid.GenerateUUID(); err != nil {
		return nil, err
	}

	err = s.withTxRetry(func(tx *gorm.DB) error {
		return tx.Create(c).Error
	})

	return c, err
}

func (s *ConversionStore) withTxRetry(fn func(tx *gorm.DB) error) error {
	return WithTxRetryDefault(fn, s.db)
}
