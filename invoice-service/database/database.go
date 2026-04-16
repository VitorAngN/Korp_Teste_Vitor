package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Importante: No teste local com SQLite, vamos usar um banco próprio para simular separação total
	// embora, se quiséssemos, a gente pudesse usar um arquivo só. Separado cumpre melhor a regra microsserviço.
	db, err := gorm.Open(sqlite.Open("korp_faturamento.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados Faturamento: %v\n", err)
	}

	DB = db
	log.Println("Conectado ao SQLite (Faturamento) com sucesso.")
}
