package caching

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/agnivade/levenshtein"
	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CachingManager struct {
	CacheDBPath string
	DB          *gorm.DB
}

func (cachingManager *CachingManager) DBInit() error {
	const driverName = "sqlite3_with_levenshtein"
	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("levenshtein", func(a, b string) int {
				return levenshtein.ComputeDistance(a, b)
			}, true)
		},
	})

	// Use that custom driver when opening the DB
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: driverName,
		DSN:        "file:" + cachingManager.CacheDBPath + "?cache=shared&mode=rwc", // Can adjust DSN as needed
	}, &gorm.Config{})

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	if db == nil {
		panic("db is nil DBInit()")
	}

	err = db.AutoMigrate(&igdb.Metadata{})
	if err != nil {
		return err
	}
	cachingManager.DB = db
	fmt.Println("Successfully connected to database: " + db.Name())
	return nil
}
func (cachingManager *CachingManager) GetMetadataListFromDBbyName(name string) []igdb.Metadata {
	db := cachingManager.DB
	var metadataList []igdb.Metadata
	var ids []uint
	db.Raw(`
    SELECT id FROM metadata
    WHERE name LIKE ? OR levenshtein(lower(name), lower(?)) <= ?`,
	"%"+name+"%", name, 5).Scan(&ids)
	if len(ids) > 0 {
		db.
			Preload("Cover").
			Preload("Artworks").
			Preload("Screenshots").
			Preload("Genres").Find(&metadataList, ids)
	}

	return metadataList
}

func (cachingManager *CachingManager) GetMetadataFromDBbyID(id int) *igdb.Metadata {
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