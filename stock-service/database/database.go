package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("korp_estoque.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados do Estoque: %v\n", err)
	}

	DB = db
	log.Println("Conectado ao SQLite (Estoque) com sucesso.")
}
