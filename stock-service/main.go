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

	// Tratamento de CORS configurado
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.GET("/api/products", handlers.GetProducts)
	r.POST("/api/products", handlers.CreateProduct)
	r.POST("/api/products/decrement", handlers.DecrementStock)

	// Requisito Opcional B: Uso de Inteligência Artificial (Mock Corporativo)
	r.POST("/api/products/ai/generate", handlers.GenerateAIDescription)

	log.Println("Serviço de Estoque rodando na porta 8081.")
	r.Run(":8081")
}
