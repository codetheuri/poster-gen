package dto

import "encoding/json"

// PosterInput is the DTO for generating a poster.
type PosterInput struct {
	BusinessName      string                 `json:"business_name" validate:"required"`
	Data              map[string]interface{} `json:"data" validate:"required"`           
	CustomizationData map[string]interface{} `json:"customization_data" validate:"omitempty"` 
}

// TemplateInput is the DTO for creating/updating a template.
type TemplateInput struct {
	Name                 string          `json:"name" validate:"required,max=100"`
	Type                 string          `json:"type" validate:"required,max=50"`
	LayoutID             uint            `json:"layout_id" validate:"required,gt=0"` 
	Price                int             `json:"price" validate:"omitempty,min=0"`
	ThumbnailURL         string          `json:"thumbnail_url" validate:"omitempty,url,max=255"`
	IsActive             bool            `json:"is_active" validate:"omitempty"`
	RequiredFields       json.RawMessage `json:"required_fields" validate:"required"`       
	DefaultCustomization json.RawMessage `json:"default_customization" validate:"required"` 
}

type AssetInput struct {
	Name         string `json:"name" validate:"required,max=100"`
	Type         string `json:"type" validate:"required,max=50"` 
	Data         string `json:"data" validate:"required"`        
	DefaultColor string `json:"default_color" validate:"omitempty,hexcolor|rgb|rgba"`
}

type LayoutInput struct {
	Name     string `json:"name" validate:"required,max=50"`
	FilePath string `json:"file_path" validate:"required,max=255"` 
}

