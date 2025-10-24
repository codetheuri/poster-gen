package repositories

import (
	"context"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)


type PosterTemplateRepository interface {
	CreateTemplate(ctx context.Context, template *models.PosterTemplate) error
	GetTemplateByID(ctx context.Context, id uint) (*models.PosterTemplate, error)
	GetActiveTemplates(ctx context.Context) ([]*models.PosterTemplate, error)
	UpdateTemplate(ctx context.Context, template *models.PosterTemplate) error
	DeleteTemplate(ctx context.Context, id uint) error
	// Add GetTemplateByName if needed
}

type posterTemplateRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewPosterTemplateRepository(db *gorm.DB, log logger.Logger) PosterTemplateRepository {
	return &posterTemplateRepository{db: db, log: log}
}

func (r *posterTemplateRepository) CreateTemplate(ctx context.Context, template *models.PosterTemplate) error {
	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		r.log.Error("Failed to create template", err, "template_name", template.Name)
		return err
	}
	return nil
}

func (r *posterTemplateRepository) GetTemplateByID(ctx context.Context, id uint) (*models.PosterTemplate, error) {
	var template models.PosterTemplate
	// Use Preload to fetch the associated Layout data automatically
	if err := r.db.WithContext(ctx).Preload("Layout").First(&template, id).Error; err != nil {
		r.log.Error("Failed to get template by ID", err, "template_id", id)
		return nil, err
	}
	return &template, nil
}

func (r *posterTemplateRepository) GetActiveTemplates(ctx context.Context) ([]*models.PosterTemplate, error) {
	var templates []*models.PosterTemplate
	// Use Preload here as well if you need Layout info in the list
	if err := r.db.WithContext(ctx).Preload("Layout").Where("is_active = ?", true).Find(&templates).Error; err != nil {
		r.log.Error("Failed to get active templates", err)
		return nil, err
	}
	return templates, nil
}

func (r *posterTemplateRepository) UpdateTemplate(ctx context.Context, template *models.PosterTemplate) error {
	if err := r.db.WithContext(ctx).Save(template).Error; err != nil {
		r.log.Error("Failed to update template", err, "template_id", template.ID)
		return err
	}
	return nil
}

func (r *posterTemplateRepository) DeleteTemplate(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.PosterTemplate{}, id).Error; err != nil {
		r.log.Error("Failed to delete template", err, "template_id", id)
		return err
	}
	return nil
}
