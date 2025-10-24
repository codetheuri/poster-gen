package dto

import "encoding/json"

// PosterResponse represents the response structure for a generated poster.
type PosterResponse struct {
	ID           uint   `json:"id"`
	TemplateID   uint   `json:"template_id"` // Corresponds to PosterTemplateID
	BusinessName string `json:"business_name"`
	PDFURL       string `json:"pdf_url"`
	Status       string `json:"status"`
}

// TemplateResponse represents the response structure for a poster template (customization profile).
type TemplateResponse struct {
	ID                   uint            `json:"id"`
	Name                 string          `json:"name"`
	Type                 string          `json:"type"`
	LayoutID             uint            `json:"layout_id"`
	LayoutFilePath       string          `json:"layout_file_path,omitempty"` // Included when fetching template
	Price                int             `json:"price"`
	ThumbnailURL         string          `json:"thumbnail_url"`
	IsActive             bool            `json:"is_active"`
	RequiredFields       json.RawMessage `json:"required_fields"`       // Send raw JSON to frontend
	DefaultCustomization json.RawMessage `json:"default_customization"` // Send raw JSON to frontend
}

// LayoutResponse represents the response structure for a layout.
type LayoutResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

// AssetResponse represents the response structure for an asset.
type AssetResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Data         string `json:"data"` // SVG or Base64
	DefaultColor string `json:"default_color,omitempty"`
}
