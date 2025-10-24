package repositories

import (
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)


type PosterRepository struct {
	LayoutRepo         LayoutRepository
	PosterTemplateRepo PosterTemplateRepository
	AssetRepo          AssetRepository
	PosterRepo         PosterSubRepository
	// OrderRepo       OrderSubRepository // Keep commented if Order model is optional
}

// NewPosterRepository constructor for the main repository aggregator.
func NewPosterRepository(db *gorm.DB, log logger.Logger) *PosterRepository {
	return &PosterRepository{
		LayoutRepo:         NewLayoutRepository(db, log),
		PosterTemplateRepo: NewPosterTemplateRepository(db, log),
		AssetRepo:          NewAssetRepository(db, log),
		PosterRepo:         NewPosterSubRepository(db, log),
		// OrderRepo:       NewOrderSubRepository(db, log), // Keep commented if Order model is optional
	}
}
