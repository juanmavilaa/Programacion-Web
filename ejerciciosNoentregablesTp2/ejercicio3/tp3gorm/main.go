package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tp3gorm/models"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=tp3_gorm port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con la base de datos")
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error al realizar la migración")
	}
}
