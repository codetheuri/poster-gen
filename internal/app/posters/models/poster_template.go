package models

import "gorm.io/gorm"

// PosterTemplate defines available designs
type PosterTemplate struct {
    gorm.Model
    Name      string `json:"name" gorm:"type:varchar(50);not null;unique"`
    Type      string `json:"type" gorm:"type:varchar(20);not null"` // e.g., "standard", "premium"
    Price     int    `json:"price" gorm:"not null;default:50"`      // In KSH
    Thumbnail string `json:"thumbnail" gorm:"type:varchar(255)"`
    IsActive  bool   `json:"is_active" gorm:"default:true"`
}