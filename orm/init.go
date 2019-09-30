package orm

import (
	//"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func NewDB(isCreateTable bool, debug bool, postgres string) error {
	db, err := gorm.Open("postgres", fmt.Sprintf("%s?sslmode=disable", postgres))
	if err != nil {
		return err
	}
	//db.DB().Ping()
	if debug {
		db.LogMode(true)
	}

	if isCreateTable {
		if !db.HasTable(&JimuTask{}) {
			db.Debug().CreateTable(&JimuTask{})
		}
		db.AutoMigrate(&JimuTask{})
	}

	DB = db

	return nil
}
