package models

import (
    "github.com/codetheuri/poster-gen/internal/app/auth/models" // Update path
    "gorm.io/gorm"
)

// Order handles payments
type Order struct {
    gorm.Model
    UserID      uint   `json:"user_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
    OrderNumber string `json:"order_number" gorm:"type:varchar(20);unique;not null"`
    TotalAmount int    `json:"total_amount" gorm:"not null"`
    Status      string `json:"status" gorm:"type:varchar(20);default:'pending'"`
    MpesaReceipt string `json:"mpesa_receipt" gorm:"type:varchar(20)"`

    // Relationships
    Posters []Poster `json:"posters" gorm:"foreignKey:OrderID"`
    User    *models.User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName customizes the table name
func (Order) TableName() string {
    return "orders"
}