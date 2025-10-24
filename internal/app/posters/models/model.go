package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	
)

type Layout struct {
	gorm.Model
	Name            string `json:"name" gorm:"type:varchar(50);not null;unique"`
	FilePath        string `json:"file_path" gorm:"type:varchar(255);not null"`
	PosterTemplates []PosterTemplate `json:"-" gorm:"foreignKey:LayoutID"`
}

func (Layout) TableName() string {
	return "layouts"
}

type PosterTemplate struct {
	gorm.Model
	Name                 string         `json:"name" gorm:"type:varchar(100);not null;unique"`
	Type                 string         `json:"type" gorm:"type:varchar(50);not null;index"`
	LayoutID             uint           `json:"layout_id" gorm:"not null;index"`
	Price                int            `json:"price" gorm:"not null;default:0"`
	ThumbnailURL         string         `json:"thumbnail_url" gorm:"type:varchar(255)"`
	IsActive             bool           `json:"is_active" gorm:"default:true;index"`
	RequiredFields       datatypes.JSON `json:"required_fields" gorm:"not null"`
	DefaultCustomization datatypes.JSON `json:"default_customization" gorm:"not null"`
	Layout               Layout         `json:"layout" gorm:"foreignKey:LayoutID"`
}

func (PosterTemplate) TableName() string {
	return "poster_templates"
}

type Asset struct {
	gorm.Model
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	Type         string `json:"type" gorm:"type:varchar(50);not null;index"`
	Data         string `json:"data" gorm:"type:text;not null"`
	DefaultColor string `json:"default_color" gorm:"type:varchar(7)"`
}

func (Asset) TableName() string {
	return "assets"
}

type Poster struct {
	gorm.Model
	PosterTemplateID   uint           `json:"poster_template_id" gorm:"not null;index"`
	BusinessName       string         `json:"business_name" gorm:"type:varchar(255);not null"`
	UserInputData      datatypes.JSON `json:"user_input_data" gorm:"not null"`
	FinalCustomization datatypes.JSON `json:"final_customization_data" gorm:"not null"`
	PDFURL             string         `json:"pdf_url" gorm:"type:varchar(255)"`
	Status             string         `json:"status" gorm:"type:varchar(50);default:'completed';index"`
	PosterTemplate     PosterTemplate `json:"poster_template" gorm:"foreignKey:PosterTemplateID"`
}

func (Poster) TableName() string {
	return "posters"
}

