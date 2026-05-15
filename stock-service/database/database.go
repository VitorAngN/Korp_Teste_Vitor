package database

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=user password=password dbname=korp port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("PostgreSQL indisponível para Estoque (%v). Usando SQLite local.", err)
		db, err = gorm.Open(sqlite.Open("korp_estoque.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Falha ao conectar no banco de dados do Estoque: %v\n", err)
		}
		log.Println("Conectado ao SQLite (Estoque) com sucesso.")
	} else {
		log.Println("Conectado ao PostgreSQL (Estoque) com sucesso.")
	}

	DB = db
}
