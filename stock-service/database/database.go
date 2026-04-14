package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := "user"
	password := "password"
	dbname := "korp_db" // We will use one DB and AutoMigrate to create tables
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados do Estoque: %v\n", err)
	}

	DB = db
	log.Println("Conectado ao PostgreSQL (Estoque) com sucesso.")
}
