package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vitor/invoice-service/database"
	"github.com/vitor/invoice-service/models"
	"github.com/vitor/invoice-service/services"
)

type CreateInvoiceRequest struct {
	Items []models.InvoiceItem `json:"items"`
}

func CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if len(req.Items) == 0 {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "A nota fiscal deve possuir ao menos 1 item."})
		return
	}

	invoice := models.Invoice{
		Status: "Aberta", // Garantir o status inicial obrigatório
		Items:  req.Items,
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar nota fiscal"})
		return
	}
	
	c.JSON(http.StatusCreated, invoice)
}

func GetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	database.DB.Preload("Items").Order("id desc").Find(&invoices)
	c.JSON(http.StatusOK, invoices)
}

// PrintInvoice lida com a impressão, atualiza o status de Aberta para Fechada e debita produtos via microsserviço de estoque
func PrintInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Começamos uma transaction para precaver concorrencia de 2 prints ao msm tempo
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var invoice models.Invoice
	if err := tx.Preload("Items").First(&invoice, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota fiscal não encontrada"})
		return
	}

	// Regras - Idempotência / Consistência
	if invoice.Status != "Aberta" {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Apenas notas com status 'Aberta' podem ser impressas."})
		return
	}

	var stockItems []services.StockItem
	for _, item := range invoice.Items {
		stockItems = append(stockItems, services.StockItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	// Como a chamada HTTP é externa, não seguramos o lock do banco absurdamente
	// mas mantemos a transação pra só fechar se tudo der certo.
	err = services.DecrementStock(stockItems)
	if err != nil {
		tx.Rollback()
		log.Printf("Falha na impressão da Nota %d. Erro Estoque: %v", invoice.Number, err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	invoice.Status = "Fechada"
	if err := tx.Save(&invoice).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gravar estado fechado."})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Nota fiscal impressa e fechada com sucesso!"})
}
