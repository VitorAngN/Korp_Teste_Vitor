package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vitor/stock-service/database"
	"github.com/vitor/stock-service/handlers"
	"github.com/vitor/stock-service/models"
)

func main() {
	log.Println("Inciando Microsserviço de Estoque...")

	// Conectar ao Banco de Dados (Postgres)
	database.ConnectDB()

	// Migrar as tabelas
	database.DB.AutoMigrate(&models.Product{})

	// Configurar rotas (Web framework: Gin)
	r := gin.Default()

	// Tratamento de CORS básico
	r.Use(cors.Default())

	r.GET("/api/products", handlers.GetProducts)
	r.POST("/api/products", handlers.CreateProduct)
	r.POST("/api/products/decrement", handlers.DecrementStock)

	log.Println("Serviço de Estoque rodando na porta 8081.")
	r.Run(":8081")
}
