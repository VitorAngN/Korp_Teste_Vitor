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
		log.Printf("PostgreSQL indisponível para Faturamento (%v). Usando SQLite local.", err)
		db, err = gorm.Open(sqlite.Open("korp_faturamento.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Falha ao conectar no banco de dados Faturamento: %v\n", err)
		}
		log.Println("Conectado ao SQLite (Faturamento) com sucesso.")
	} else {
		log.Println("Conectado ao PostgreSQL (Faturamento) com sucesso.")
	}

	DB = db
}
