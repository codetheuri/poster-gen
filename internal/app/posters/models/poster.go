package models

import (
	"github.com/codetheuri/poster-gen/internal/app/auth/models" // Update path
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Poster represents a generated payment poster
// type Poster struct {
//     gorm.Model
//     UserID     uint           `json:"user_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
//     OrderID    uint           `json:"order_id" gorm:"not null;index"`
//     TemplateID uint           `json:"template_id" gorm:"not null;index"`
//     BusinessName string       `json:"business_name" gorm:"type:varchar(100);not null"`
//     PhoneNumber  string       `json:"phone_number" gorm:"type:varchar(15)"` // Optional override
//     PaymentType   string       `json:"payment_type" gorm:"type:varchar(20);not null"` // e.g., "mpesa", "bank"
//     TillNumber    *string      `json:"till_number" gorm:"type:varchar(10)"`           // Nullable
//     PaybillNumber *string      `json:"paybill_number" gorm:"type:varchar(10)"`        // Nullable
//     AccountNumber *string      `json:"account_number" gorm:"type:varchar(20)"`        // Nullable
//     AccountName   *string      `json:"account_name" gorm:"type:varchar(100)"`         // Nullable
//     BankName      *string      `json:"bank_name" gorm:"type:varchar(50)"`             // Nullable
//     PrimaryColor   string      `json:"primary_color" gorm:"type:varchar(7);default:'#0369a1'"`
//     SecondaryColor string      `json:"secondary_color" gorm:"type:varchar(7);default:'#0ea5e9'"`
//     LogoURL        string      `json:"logo_url" gorm:"type:varchar(255)"`
//     QRType         string      `json:"qr_type" gorm:"type:varchar(20);default:'payment'"`
//     QRData         string      `json:"qr_data" gorm:"type:varchar(255)"`
//     PDFURL         string      `json:"pdf_url" gorm:"type:varchar(255)"`
//     ImageURL       string      `json:"image_url" gorm:"type:varchar(255)"`
//     Status         string      `json:"status" gorm:"type:varchar(20);default:'draft'"`

//     // Relationships
//     Template PosterTemplate `json:"template" gorm:"foreignKey:TemplateID"`
//     User     *models.User   `json:"user" gorm:"foreignKey:UserID"`
//     Order    *Order         `json:"order" gorm:"foreignKey:OrderID"`
// }

type Poster struct {
	gorm.Model

	UserID     *uint  `json:"user_id,omitempty" gorm:"index;constraint:OnDelete:SET NULL"`
    OrderID    *uint  `json:"order_id,omitempty" gorm:"index"`
	TemplateID   uint   `json:"template_id" gorm:"not null;index"`
	BusinessName string `json:"business_name" gorm:"type:varchar(100);not null"` // Keep common fields
	DynamicData datatypes.JSON `json:"dynamic_data" gorm:"not null"`
	PDFURL string `json:"pdf_url" gorm:"type:varchar(255)"`
	Status string `json:"status" gorm:"type:varchar(20);default:'draft'"`

	// Relationships
	Template PosterTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	User     *models.User   `json:"user" gorm:"foreignKey:UserID"`
	Order    *Order         `json:"order" gorm:"foreignKey:OrderID"`
}

func (Poster) TableName() string {
	return "posters"
}
