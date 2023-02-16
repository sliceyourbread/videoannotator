package controllers

import (
	"fmt"
	"log"
	"os"
	"papercup/videoannotator/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	Conn *gorm.DB
}

var DB Server

func Connect(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("Connected to DB")
	db.AutoMigrate(&models.Video{}, &models.User{}, &models.Annotation{})

	DB = Server{
		Conn: db,
	}
}
