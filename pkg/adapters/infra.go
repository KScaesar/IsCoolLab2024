package adapters

import (
	"gorm.io/gorm"
)

type Infra struct {
	Database *gorm.DB
}

func (infra *Infra) Cleanup() {
	db, _ := infra.Database.DB()
	err := db.Close()
	if err != nil {
		panic(err)
	}
}
