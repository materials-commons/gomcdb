package mcmodel

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestQueryDataset(t *testing.T) {
	dsn := os.Getenv("MC_DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open db: %s", err)
	}

	var ds Dataset
	result := db.Find(&ds, 1)

	require.NoError(t, result.Error, "Query returned error: %s", result.Error)
	fmt.Printf("%+v\n", ds)
	fs, err := ds.GetFileSelection()
	require.NoError(t, err, "GetFileSelection returned error: %s", err)
	fmt.Printf("%+v\n", fs)
}

func TestBuildingEntitiesQuery(t *testing.T) {
	dsn := "mc:mcpw@tcp(127.0.0.1:3306)/mc?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := os.Getenv("MC_DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open db: %s", err)
	}

	d := Dataset{ID: 1}

	experimentIdsSubquery := db.Table("item2entity_selection").
		Select("experiment_id").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset")

	entityIdsFromExperimentSubquery := db.Table("experiment2entity").
		Select("entity_id").
		Where("experiment_id in (?)", experimentIdsSubquery)

	entityNamesFromExperimentSubquery := db.Table("item2entity_selection").
		Select("entity_name").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset").
		Where("experiment_id in (?)", experimentIdsSubquery)

	entityIdSubquery := db.Table("item2entity_selection").
		Select("entity_id").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset")

	var entities []Entity
	stmt := db.Preload("Files.Directory").
		Where("id in (?)", entityIdsFromExperimentSubquery).
		Where("name in (?)", entityNamesFromExperimentSubquery).
		Or("id in (?)", entityIdSubquery).
		Find(&entities).Statement
	fmt.Println(stmt.SQL.String())
}

func TestDataset_GetEntitiesFromTemplate(t *testing.T) {
	dsn := "mc:mcpw@tcp(127.0.0.1:3306)/mc?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := os.Getenv("MC_DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open db: %s", err)
	}

	var ds Dataset
	result := db.Find(&ds, 6)
	require.NoError(t, result.Error, "Query returned error: %s", result.Error)
	entities, err := ds.GetEntitiesFromTemplate(db)
	require.NoError(t, err, "GetEntitiesFromTemplate failed: %s\n", err)
	require.NotEmpty(t, entities, "Entities is empty")
	for _, entity := range entities {
		if entity.ID == 2324 {
			require.NotEmpty(t, entity.Files, "entity %d has empty files", entity.ID)
		}
	}
}

func TestDatasetTime(t *testing.T) {
	dsn := "mc:mcpw@tcp(127.0.0.1:3306)/mc?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := os.Getenv("MC_DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open db: %s", err)
	}

	var ds Dataset
	result := db.Find(&ds, 3)
	require.NoError(t, result.Error, "Query returned error: %s", result.Error)
	fmt.Printf("%+v\n", ds)

	fmt.Println("time is: ", ds.PublishedAt.IsZero())
}
