package repositories

import (
	"context"

	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)

// PosterRepository defines the interface for poster data operations.
// Renamed from PosterSubRepository.
type PosterSubRepository interface {
	CreatePoster(ctx context.Context, poster *models.Poster) error
	GetPosterByID(ctx context.Context, id uint) (*models.Poster, error)
	// Add other methods as needed (Update, Delete, ListByUser, etc.)
}

type posterRepository struct {
	db  *gorm.DB
	log logger.Logger
}

// NewPosterRepository creates a new PosterRepository.
// Renamed from NewPosterSubRepository.
func NewPosterSubRepository(db *gorm.DB, log logger.Logger) PosterSubRepository {
	return &posterRepository{db: db, log: log}
}

func (r *posterRepository) CreatePoster(ctx context.Context, poster *models.Poster) error {
	if err := r.db.WithContext(ctx).Create(poster).Error; err != nil {
		r.log.Error("Failed to create poster", err)
		return err
	}
	return nil
}

func (r *posterRepository) GetPosterByID(ctx context.Context, id uint) (*models.Poster, error) {
	var poster models.Poster
	// Preload related data if needed when retrieving a specific poster
	if err := r.db.WithContext(ctx).Preload("PosterTemplate.Layout").First(&poster, id).Error; err != nil {
		r.log.Error("Failed to get poster by ID", err, "poster_id", id)
		return nil, err
	}
	return &poster, nil
}

// Add UpdatePoster, DeletePoster implementations if needed
