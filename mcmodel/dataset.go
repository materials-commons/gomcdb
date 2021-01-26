package mcmodel

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Dataset struct {
	ID            int
	UUID          string
	ProjectID     int
	License       string
	LicenseLink   string
	Description   string
	Summary       string
	DOI           string
	Authors       string
	FileSelection string
}

//func (f *FileSelection) Scan(src interface{}) error {
//	fmt.Println("Calling Scan")
//	// The data stored in a JSON field is actually returned as []uint8
//	val := src.([]byte)
//	var fs FileSelection
//	err := json.Unmarshal(val, &fs)
//	if err != nil {
//		fmt.Println("json.Unmarshal failed:", err)
//		return err
//	}
//	*f = fs
//	return nil
//}

func (d Dataset) GetFileSelection() (*FileSelection, error) {
	var fs FileSelection
	err := json.Unmarshal([]byte(d.FileSelection), &fs)
	return &fs, err
}

type FileSelection struct {
	IncludeFiles []string `json:"include_files"`
	ExcludeFiles []string `json:"exlude_files"`
	IncludeDirs  []string `json:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs"`
}

func (d Dataset) GetEntitiesFromTemplate(db *gorm.DB) ([]Entity, error) {
	experimentIdsSubSubquery := db.Table("item2entity_selection").
		Select("experiment_id").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset")

	entityIdsFromExperimentSubquery := db.Table("experiment2entity").
		Select("entity_id").
		Where("experiment_id in (?)", experimentIdsSubSubquery)

	entityNamesFromExperimentSubquery := db.Table("item2entity_selection").
		Select("entity_name").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset").
		Where("experiment_id in (?)", experimentIdsSubSubquery)

	entityIdSubquery := db.Table("item2entity_selection").
		Select("entity_id").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset")

	var entities []Entity
	result := db.Preload("Files.Directory").
		Where("id in (?)", entityIdsFromExperimentSubquery).
		Where("name in (?)", entityNamesFromExperimentSubquery).
		Or("id in (?)", entityIdSubquery).
		Find(&entities)

	return entities, result.Error
}
