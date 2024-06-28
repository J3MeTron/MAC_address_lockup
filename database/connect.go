package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"myapp/models"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open("mac_addresses.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Миграция схемы
	db.AutoMigrate(&models.Device{})

	DB = db
}
