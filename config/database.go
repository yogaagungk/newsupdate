package config

import (
	"log"

	"github.com/yogaagungk/newsupdate/model"

	"github.com/jinzhu/gorm"
)

// OpenDB , konfigurasi dan open connection ke database,
// Disini database yang digunakan adalah mysql
func OpenDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/newsupdate?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Panic("Error while connecting to database")
	}

	db.DB().SetMaxIdleConns(10)
	db.AutoMigrate(&model.News{})

	return db
}
