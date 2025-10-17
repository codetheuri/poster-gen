package repositories

import (
	"context"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)

// PosterTemplateSubRepository interface for template operations
type PosterTemplateSubRepository interface {
	CreateTemplate(ctx context.Context, template *models.PosterTemplate) error
	GetTemplateByID(ctx context.Context, id uint) (*models.PosterTemplate, error)
	GetActiveTemplates(ctx context.Context) ([]*models.PosterTemplate, error)
	UpdateTemplate(ctx context.Context, template *models.PosterTemplate) error
	DeleteTemplate(ctx context.Context, id uint) error
}


type posterTemplateSubRepository struct {
	db  *gorm.DB
	log logger.Logger
}

// NewPosterTemplateSubRepository constructor
func NewPosterTemplateSubRepository(db *gorm.DB, log logger.Logger) PosterTemplateSubRepository {
	return &posterTemplateSubRepository{
		db:  db,
		log: log,
	}
}

func (r *posterTemplateSubRepository) CreateTemplate(ctx context.Context, template *models.PosterTemplate) error {
	r.log.Info("Creating poster template", "name", template.Name)
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *posterTemplateSubRepository) GetTemplateByID(ctx context.Context, id uint) (*models.PosterTemplate, error) {
	r.log.Info("Getting poster template by ID", "id", id)
	var template models.PosterTemplate
	err := r.db.WithContext(ctx).First(&template, id).Error
	if err != nil {
		r.log.Error("Failed to get template by ID", err, "id", id)
		return nil, err
	}
	return &template, nil
}

func (r *posterTemplateSubRepository) GetActiveTemplates(ctx context.Context) ([]*models.PosterTemplate, error) {
	r.log.Info("Getting active poster templates")
	var templates []*models.PosterTemplate
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&templates).Error
	if err != nil {
		r.log.Error("Failed to get active templates", err)
		return nil, err
	}
	return templates, nil
}

func (r *posterTemplateSubRepository) UpdateTemplate(ctx context.Context, template *models.PosterTemplate) error {
	r.log.Info("Updating poster template", "id", template.ID)
	// return r.db.WithContext(ctx).Save(template).Error
       return r.db.WithContext(ctx).Model(&models.PosterTemplate{}).Where("id = ?", template.ID).Updates(template).Error	
// 
}

func (r *posterTemplateSubRepository) DeleteTemplate(ctx context.Context, id uint) error {
	r.log.Info("Deleting poster template", "id", id)
	return r.db.WithContext(ctx).Delete(&models.PosterTemplate{}, id).Error
}
