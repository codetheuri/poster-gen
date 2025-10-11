package dto

// PosterResponse represents the response structure for a poster
type PosterResponse struct {
    ID          uint   `json:"id"`
    UserID      uint   `json:"user_id"`
    OrderID     uint   `json:"order_id"`
    TemplateID  uint   `json:"template_id"`
    BusinessName string `json:"business_name"`
    PDFURL      string `json:"pdf_url"`
    Status      string `json:"status"`
}

// TemplateResponse represents the response structure for a template
type TemplateResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Type      string `json:"type"`
    Price     int    `json:"price"`
    Thumbnail string `json:"thumbnail"`
    IsActive  bool   `json:"is_active"`
}

// OrderResponse represents the response structure for an order
type OrderResponse struct {
    ID          uint   `json:"id"`
    UserID      uint   `json:"user_id"`
    OrderNumber string `json:"order_number"`
    TotalAmount int    `json:"total_amount"`
    Status      string `json:"status"`
    MpesaReceipt *string `json:"mpesa_receipt,omitempty"`
}