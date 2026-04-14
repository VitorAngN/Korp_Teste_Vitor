package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type StockItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

type DecrementRequest struct {
	Items []StockItem `json:"items"`
}

// DecrementStock sends an HTTP request to Stock microservice with Retry resiliency
func DecrementStock(items []StockItem) error {
	stockServiceURL := os.Getenv("STOCK_SERVICE_URL")
	if stockServiceURL == "" {
		stockServiceURL = "http://localhost:8081"
	}
	url := fmt.Sprintf("%s/api/products/decrement", stockServiceURL)

	reqBody := DecrementRequest{Items: items}
	jsonData, _ := json.Marshal(reqBody)

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			
			// Se o microsserviço retornar regra de negocio quebrada (sem saldo, etc), nós não damos retry.
			if resp.StatusCode == http.StatusNotFound {
			    return errors.New("Produto não encontrado no estoque")
			}
			if resp.StatusCode == http.StatusConflict {
				return errors.New("Saldo insuficiente no estoque")
			}
			
			// Para outros erros 500, tentamos denovo
			fmt.Printf("Serviço de estoque respondeu %d. Retentando...\n", resp.StatusCode)
		} else {
			fmt.Printf("Falha ao contatar serviço de Estoque: %v. Retentando em 2 segundos...\n", err)
		}
		
		time.Sleep(2 * time.Second)
	}
	
	// Feedback de falha apropriado ao usuário (Microserviços)
	return errors.New("O serviço de controle de Estoque está indisponível no momento. A impressão da sua nota falhou. Tente novamente mais tarde.")
}
