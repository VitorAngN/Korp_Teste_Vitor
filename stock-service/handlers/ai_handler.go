package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AIRequest struct {
	ProductName string `json:"product_name"`
}

// GenerateAIDescription simula a chamada a um modelo de IA (LLM) e devolve uma descrição baseada no contexto industrial do Korp Teste.
func GenerateAIDescription(c *gin.Context) {
	var req AIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required for AI generation."})
		return
	}

	// Simula a latência de rede/processamento de uma API externa (OpenAI / Gemini)
	time.Sleep(1200 * time.Millisecond)

	nameUpper := strings.ToUpper(req.ProductName)

	// Contextos randômicos corporativos para dar ar de inteligência ao Mock:
	benefits := []string{
		"alta durabilidade e performance otimizada para o setor industrial",
		"certificação de qualidade internacional e baixo custo de manutenção",
		"design ergonômico com tecnologia de ponta e eficiência energética",
		"integração contínua com maquinário moderno e alto rendimento operacional",
	}

	rand.Seed(time.Now().UnixNano())
	chosenBenefit := benefits[rand.Intn(len(benefits))]

	generatedText := fmt.Sprintf("O produto %s é desenvolvido com %s, garantindo aderência absoluta aos padrões de excelência Korp.", nameUpper, chosenBenefit)

	c.JSON(http.StatusOK, gin.H{
		"description": generatedText,
	})
}
