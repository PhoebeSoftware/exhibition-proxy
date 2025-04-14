package db_models

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ImageID string `json:"image_id"`
}

type Genre struct {
	GenreID int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
}

type DBMetadata struct {
	Id          int `gorm:"primaryKey"`
	IGDBID      int
	Name        string
	Description string
	CoverID     uint
	Cover       Image   `gorm:"foreignKey:CoverID"`
	Artworks    []Image `gorm:"many2many:db_metadata_artworks"`
	Screenshots []Image `gorm:"many2many:db_metadata_screenshots"`
	Genres      []Genre `gorm:"many2many:db_metadata_genres"`
}
