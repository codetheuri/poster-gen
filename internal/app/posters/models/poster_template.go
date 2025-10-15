package models

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

// PosterTemplate defines available designs
// type PosterTemplate struct {
//     gorm.Model
//     Name      string `json:"name" gorm:"type:varchar(50);not null;unique"`
//     Type      string `json:"type" gorm:"type:varchar(20);not null"` // e.g., "standard", "premium"
//     Price     int    `json:"price" gorm:"not null;default:50"`      // In KSH
//     Thumbnail string `json:"thumbnail" gorm:"type:varchar(255)"`
//     IsActive  bool   `json:"is_active" gorm:"default:true"`
//     Layout    string `json:"layout" gorm:"not null"` // JSON or XML layout definition

// }
type PosterTemplate struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(50);not null;unique"`
	Type      string `json:"type" gorm:"type:varchar(20);not null"`
	Price     int    `json:"price" gorm:"not null;default:50"`
	Thumbnail string `json:"thumbnail" gorm:"type:varchar(255)"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	Layout    string `json:"layout" gorm:"not null"`
	// --- NEW FIELD ---
	// This will define the dynamic form fields needed for this template.
	RequiredFields datatypes.JSON `json:"required_fields" gorm:"not null"`
}
