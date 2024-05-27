package pkg

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteGorm(dsn string) (*gorm.DB, error) {
	// config := gorm.Config{}
	config := gorm.Config{Logger: logger.Discard}
	return gorm.Open(sqlite.Open(dsn), &config)
}
