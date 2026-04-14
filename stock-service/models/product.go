package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	Description string    `gorm:"not null" json:"description"`
	Balance     int       `gorm:"not null;default:0" json:"balance"` // Saldo em estoque
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
