package repositories

import (
    "context"
    "github.com/codetheuri/poster-gen/internal/app/posters/models"
    "github.com/codetheuri/poster-gen/pkg/logger"
    "gorm.io/gorm"
)

// PosterSubRepository interface for poster operations
type PosterSubRepository interface {
    CreatePoster(ctx context.Context, poster *models.Poster) error
    GetPosterByID(ctx context.Context, id uint) (*models.Poster, error)
    UpdatePoster(ctx context.Context, poster *models.Poster) error
    DeletePoster(ctx context.Context, id uint) error
    // Add more as needed, e.g., GetPostersByUserID, GetPostersByOrderID
}

type posterSubRepository struct {
    db  *gorm.DB
    log logger.Logger
}

// NewPosterSubRepository constructor
func NewPosterSubRepository(db *gorm.DB, log logger.Logger) PosterSubRepository {
    return &posterSubRepository{
        db:  db,
        log: log,
    }
}

func (r *posterSubRepository) CreatePoster(ctx context.Context, poster *models.Poster) error {
    r.log.Info("Creating poster", "business_name", poster.BusinessName, "user_id", poster.UserID)
    return r.db.WithContext(ctx).Create(poster).Error
}

func (r *posterSubRepository) GetPosterByID(ctx context.Context, id uint) (*models.Poster, error) {
    r.log.Info("Getting poster by ID", "id", id)
    var poster models.Poster
    err := r.db.WithContext(ctx).First(&poster, id).Error
    if err != nil {
        r.log.Error("Failed to get poster by ID", err, "id", id)
        return nil, err
    }
    return &poster, nil
}

func (r *posterSubRepository) UpdatePoster(ctx context.Context, poster *models.Poster) error {
    r.log.Info("Updating poster", "id", poster.ID)
    return r.db.WithContext(ctx).Save(poster).Error
}

func (r *posterSubRepository) DeletePoster(ctx context.Context, id uint) error {
    r.log.Info("Deleting poster", "id", id)
    return r.db.WithContext(ctx).Delete(&models.Poster{}, id).Error
}