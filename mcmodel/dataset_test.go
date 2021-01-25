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
