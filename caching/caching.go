package caching

import (
	"errors"
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

func (cachingManager *CachingManager) GetMetadataFromDB(id int) *igdb.Metadata {
	db := cachingManager.DB
	var metadata igdb.Metadata
	result := db.
		Preload("Cover").
		Preload("Artworks").
		Preload("Screenshots").
		Preload("Genres").
		First(&metadata, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &metadata
}

func (cachingManager *CachingManager) AddMetadataToDB(metadata *igdb.Metadata) {
	db := cachingManager.DB
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(metadata)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}