package config

import (
	"fmt"
	"h8-p2-finalproj-app/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDBInstance() *gorm.DB {
	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Car{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Rental{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Payment{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.TopUp{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
