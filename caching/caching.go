package caching

import (
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CachingManager struct {
	DB *gorm.DB
}

func (cachingManager *CachingManager) DBInit() error {
	db, err := gorm.Open(sqlite.Open("./cache.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Error opening database: %w", err)
	}

	err = db.AutoMigrate(&igdb.Metadata{})
	if err != nil {
		return err
	}
	cachingManager.DB = db
	fmt.Println("Successfully connected to database: " + db.Name())
	return nil
}

func (cachingManager *CachingManager) AddMetadataToDatabase(metadata *igdb.Metadata) {
	db := cachingManager.DB
	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(metadata)
}