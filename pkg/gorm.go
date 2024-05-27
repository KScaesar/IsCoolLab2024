package pkg

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteGorm(dsn string, debug bool) (*gorm.DB, error) {
	if !debug {
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Discard})
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db.Debug(), nil
}
