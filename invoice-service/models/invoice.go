package models

import "time"

// Invoice model for Faturamento
type Invoice struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Number    int           `gorm:"uniqueIndex" json:"number"`
	Status    string        `gorm:"not null;default:'Aberta'" json:"status"`
	Items     []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type InvoiceItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	InvoiceID   uint      `gorm:"not null" json:"-"`
	ProductCode string    `gorm:"not null" json:"product_code"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
