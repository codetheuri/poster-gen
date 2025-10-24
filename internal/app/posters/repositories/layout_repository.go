package repositories

import (
	"context"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)


type LayoutRepository interface {
		CreateLayout(ctx context.Context, layout *models.Layout) error 
	ListLayouts(ctx context.Context) ([]*models.Layout, error) 
	GetLayoutByID(ctx context.Context, id uint) (*models.Layout, error)
	GetLayoutByName(ctx context.Context, name string) (*models.Layout, error)
	
}

type layoutRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewLayoutRepository(db *gorm.DB, log logger.Logger) LayoutRepository {
	return &layoutRepository{db: db, log: log}
}
func (r *layoutRepository) CreateLayout(ctx context.Context, layout *models.Layout) error {
	if err := r.db.WithContext(ctx).Create(layout).Error; err != nil {
		r.log.Error("Failed to create layout", err, "layout_name", layout.Name)
		return err // Let service layer wrap the error
	}
	return nil
}
func (r *layoutRepository) ListLayouts(ctx context.Context) ([]*models.Layout, error) {
	var layouts []*models.Layout
	if err := r.db.WithContext(ctx).Find(&layouts).Error; err != nil {
		r.log.Error("Failed to list layouts", err)
		return nil, err
	}
	return layouts, nil
}


func (r *layoutRepository) GetLayoutByID(ctx context.Context, id uint) (*models.Layout, error) {
	var layout models.Layout
	if err := r.db.WithContext(ctx).First(&layout, id).Error; err != nil {
		r.log.Error("Failed to get layout by ID", err, "layout_id", id)
		return nil, err 
	}
	return &layout, nil
}

func (r *layoutRepository) GetLayoutByName(ctx context.Context, name string) (*models.Layout, error) {
	var layout models.Layout
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&layout).Error; err != nil {
		r.log.Error("Failed to get layout by name", err, "layout_name", name)
		return nil, err
	}
	return &layout, nil
}
