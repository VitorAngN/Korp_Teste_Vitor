package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitor/stock-service/database"
	"github.com/vitor/stock-service/models"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Erro ao criar produto: %v", err)
		return
	}
	c.JSON(http.StatusCreated, product)
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Order("id asc").Find(&products)
	c.JSON(http.StatusOK, products)
}

type DecrementRequest struct {
	Items []struct {
		ProductCode string `json:"product_code"`
		Quantity    int    `json:"quantity"`
	} `json:"items"`
}

// DecrementStock abate os itens do saldo garantindo atomicidade na concorrência
func DecrementStock(c *gin.Context) {
	var req DecrementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range req.Items {
		// Validar existência do produto
		var count int64
		tx.Model(&models.Product{}).Where("code = ?", item.ProductCode).Count(&count)
		if count == 0 {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado: " + item.ProductCode})
			return
		}

		// Abatimento atômico: Trata falhas de concorrência garantindo que o saldo seja >= quantidade exigida
		result := tx.Model(&models.Product{}).
			Where("code = ? AND balance >= ?", item.ProductCode, item.Quantity).
			Update("balance", gorm.Expr("balance - ?", item.Quantity))

		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao atualizar saldo"})
			return
		}

		// Se linhas afetadas for 0, mas sabemos que ele existe, significa que o saldo era insuficiente (balance < quantity)
		if result.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Saldo insuficiente para o produto: " + item.ProductCode})
			return
		}

		log.Printf("Estoque de %s debitado em %d unidade(s)", item.ProductCode, item.Quantity)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Estoque atualizado com sucesso"})
}
