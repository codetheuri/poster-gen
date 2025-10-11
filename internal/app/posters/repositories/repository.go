package repositories

import (
    "github.com/codetheuri/poster-gen/pkg/logger"
    "gorm.io/gorm"
)

// PosterRepository is the main repository aggregator for the posters module
type PosterRepository struct {
    PosterRepo          PosterSubRepository
    PosterTemplateRepo  PosterTemplateSubRepository
    OrderRepo           OrderSubRepository
}

// NewPosterRepository constructor for the main repository
func NewPosterRepository(db *gorm.DB, log logger.Logger) *PosterRepository {
    return &PosterRepository{
        PosterRepo:         NewPosterSubRepository(db, log),
        PosterTemplateRepo: NewPosterTemplateSubRepository(db, log),
        OrderRepo:          NewOrderSubRepository(db, log),
    }
}