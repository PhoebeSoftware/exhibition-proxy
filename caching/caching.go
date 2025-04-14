package caching

import (
	"exhibtion-proxy/db_models"
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

	err = db.AutoMigrate(&db_models.DBMetadata{})
	if err != nil {
		return err
	}
	cachingManager.DB = db
	fmt.Println("Succesfully connected to databse: " + db.Name())
	return nil
}

func (cachingManager *CachingManager) AddMetadataToDatabase(metadata igdb.Metadata) {
	db := cachingManager.DB

	dbMetadata := MetadataToDBMetadata(metadata)
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "igdb_id"}},
		DoNothing: true,
	}).Create(&dbMetadata)

}

func ConvertToApiGame(dbGame *db_models.DBMetadata) {

}

func MetadataToDBMetadata(metadata igdb.Metadata) *db_models.DBMetadata {
	dbGame := &db_models.DBMetadata{
		IGDBID:      metadata.Id,
		Name:        metadata.Name,
		Description: metadata.Description,
		Cover:       db_models.Image{ImageID: metadata.Cover.ImageID},
		Artworks:    ConvertIGDBImageToDBImage(metadata.Artworks),
		Screenshots: ConvertIGDBImageToDBImage(metadata.Screenshots),
		Genres: ConvertIGDBGenreToDBGenre(metadata.Genres),
	}
	return dbGame
}
func ConvertIGDBGenreToDBGenre(genreList []igdb.Genre) []db_models.Genre {
	var result []db_models.Genre
	for _, genre := range genreList {
		result = append(result, db_models.Genre{
			GenreID: genre.GenreID,
			Name: genre.Name,
		})
	}
	return result
}

func ConvertIGDBImageToDBImage(imageList []igdb.Image) []db_models.Image {
	var result []db_models.Image
	for _, image := range imageList {
		result = append(result, db_models.Image{
			ImageID: image.ImageID,
		})
	}
	return result
}
