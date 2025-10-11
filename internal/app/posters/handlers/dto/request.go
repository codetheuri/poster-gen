package dto

// PosterInput matches the service's PosterInput for request validation
type PosterInput struct {
	BusinessName  string `json:"business_name" validate:"required,max=100"`
	PhoneNumber   string `json:"phone_number" validate:"omitempty,max=15"`
	PaymentType   string `json:"payment_type" validate:"required,oneof=mpesa bank"`
	TillNumber    string `json:"till_number" validate:"omitempty,max=10"`
	PaybillNumber string `json:"paybill_number" validate:"omitempty,max=10"`
	AccountNumber string `json:"account_number" validate:"omitempty,max=20"`
	AccountName   string `json:"account_name" validate:"omitempty,max=100"`
	BankName      string `json:"bank_name" validate:"omitempty,max=50"`
	LogoURL       string `json:"logo_url" validate:"omitempty,url,max=255"`
}

// TemplateInput matches the service's TemplateInput for request validation
type TemplateInput struct {
	Name      string `json:"name" validate:"required,max=50"`
	Type      string `json:"type" validate:"required,max=20"`
	Price     int    `json:"price" validate:"required,min=0"`
	Thumbnail string `json:"thumbnail" validate:"omitempty,url,max=255"`
	IsActive  bool   `json:"is_active" validate:"omitempty"`
}

// OrderInput matches the service's OrderInput for request validation
type OrderInput struct {
	TotalAmount int `json:"total_amount" validate:"required,min=0"`
}
