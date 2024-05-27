package database

import (
	"gorm.io/gorm"

	"github.com/KScaesar/IsCoolLab2024/pkg"
	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

type GormConfing struct {
	Dsn     string
	Migrate bool
	Debug   bool
}

func NewGrom(conf *GormConfing) (*gorm.DB, error) {
	db, err := pkg.NewSqliteGorm(conf.Dsn, conf.Debug)
	if err != nil {
		return nil, err
	}

	if !conf.Migrate {
		return db, nil
	}

	err = db.AutoMigrate(
		app.User{},
		app.FileSystem{},
		app.Folder{},
		app.File{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
