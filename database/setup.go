package database

import (
	"os"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	globalInstance *gorm.DB
)

func Setup() (*gorm.DB, error) {

	if globalInstance != nil {
		return globalInstance, nil
	}

	err_load := godotenv.Load()

	if err_load != nil {
		println("Error to load the .env file")
		return nil, err_load
	}

	connection := os.Getenv("DB_CONNECTION")

	db, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		println("Error connecting to the database")
		return nil, err
	}

	db.AutoMigrate(&models.User{}, &models.Photo{})

	globalInstance = db
	return db, nil
}
