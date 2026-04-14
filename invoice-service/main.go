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

	// Tratamento de CORS básico
	r.Use(cors.Default())

	r.GET("/api/invoices", handlers.GetInvoices)
	r.POST("/api/invoices", handlers.CreateInvoice)
	r.POST("/api/invoices/:id/print", handlers.PrintInvoice)

	log.Println("Serviço de Faturamento rodando na porta 8082.")
	r.Run(":8082")
}
