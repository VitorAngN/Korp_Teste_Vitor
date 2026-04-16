package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vitor/invoice-service/database"
	"github.com/vitor/invoice-service/handlers"
	"github.com/vitor/invoice-service/models"
)

func main() {
	log.Println("Inciando Microsserviço de Faturamento...")

	database.ConnectDB()

	database.DB.AutoMigrate(&models.Invoice{}, &models.InvoiceItem{})

	r := gin.Default()

	// Tratamento de CORS configurado para aceitar headers customizados (X-Simulate-Failure)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Simulate-Failure"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.GET("/api/invoices", handlers.GetInvoices)
	r.POST("/api/invoices", handlers.CreateInvoice)
	r.POST("/api/invoices/:id/print", handlers.PrintInvoice)

	log.Println("Serviço de Faturamento rodando na porta 8082.")
	r.Run(":8082")
}
