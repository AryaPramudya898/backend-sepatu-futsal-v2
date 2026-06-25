package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID        uint              `gorm:"not null;index" json:"user_id"`
	Reference     string            `gorm:"type:varchar(100);not null;uniqueIndex" json:"reference"`
	Amount        float64           `gorm:"not null" json:"amount"`
	Description   string            `gorm:"type:text" json:"description"`
	Status        string            `gorm:"type:varchar(50);not null;default:'pending'" json:"status"` // 'pending', 'success', 'failed', 'cancelled'
	TransactionID string            `gorm:"type:varchar(100)" json:"transaction_id"`
	Items         []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
}

type TransactionItem struct {
	gorm.Model
	TransactionID uint    `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint    `gorm:"not null;index" json:"product_id"`
	ProductName   string  `gorm:"type:varchar(255);not null" json:"product_name"`
	Price         float64 `gorm:"not null" json:"price"`
	Quantity      int     `gorm:"not null" json:"quantity"`
}

type CreateTransactionRequest struct {
	Reference   string  `json:"reference" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

type UpdateTransactionStatusRequest struct {
	Status        string `json:"status" binding:"required"`
	TransactionID string `json:"transaction_id"`
}
