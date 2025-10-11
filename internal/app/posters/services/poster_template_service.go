package services

import (
	"context"
	stdErrors "errors" // Import standard errors package

	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/gorm"
)

// PosterTemplateSubService interface for template operations
type PosterTemplateSubService interface {
	CreateTemplate(ctx context.Context, input *TemplateInput) (*models.PosterTemplate, error)
	GetTemplateByID(ctx context.Context, id uint) (*models.PosterTemplate, error)
	GetActiveTemplates(ctx context.Context) ([]*models.PosterTemplate, error)
	UpdateTemplate(ctx context.Context, template *models.PosterTemplate) error
	DeleteTemplate(ctx context.Context, id uint) error
}

// TemplateInput defines the input data for creating a template
type TemplateInput struct {
	Name      string `json:"name" validate:"required,max=50"`
	Type      string `json:"type" validate:"required,max=20"`
	Price     int    `json:"price" validate:"required,min=0"`
	Thumbnail string `json:"thumbnail" validate:"omitempty,url,max=255"`
	IsActive  bool   `json:"is_active" validate:"omitempty"`
}

type posterTemplateSubService struct {
	repo      repositories.PosterTemplateSubRepository
	validator *validators.Validator
	log       logger.Logger
}

// NewPosterTemplateSubService constructor
func NewPosterTemplateSubService(repo repositories.PosterTemplateSubRepository, validator *validators.Validator, log logger.Logger) PosterTemplateSubService {
	return &posterTemplateSubService{
		repo:      repo,
		validator: validator,
		log:       log,
	}
}

func (s *posterTemplateSubService) CreateTemplate(ctx context.Context, input *TemplateInput) (*models.PosterTemplate, error) {
	s.log.Info("Creating poster template", "name", input.Name)

	validationErrors := s.validator.Struct(input)
	if validationErrors != nil {
		s.log.Warn("Validation failed for template input", validationErrors)
		return nil, errors.ValidationError("invalid template input", nil, validationErrors)
	}

	template := &models.PosterTemplate{
		Name:      input.Name,
		Type:      input.Type,
		Price:     input.Price,
		Thumbnail: input.Thumbnail,
		IsActive:  input.IsActive,
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		s.log.Error("Failed to save template to database", err)
		return nil, errors.DatabaseError("failed to save template", err)
	}

	s.log.Info("Template created successfully", "template_id", template.ID)
	return template, nil
}

func (s *posterTemplateSubService) GetTemplateByID(ctx context.Context, id uint) (*models.PosterTemplate, error) {
	s.log.Info("Getting template by ID", "id", id)
	template, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Template not found", "id", id)
			return nil, errors.NotFoundError("template not found", err)
		}
		s.log.Error("Failed to get template", err, "id", id)
		return nil, errors.DatabaseError("failed to retrieve template", err)
	}
	return template, nil
}

func (s *posterTemplateSubService) GetActiveTemplates(ctx context.Context) ([]*models.PosterTemplate, error) {
	s.log.Info("Getting active templates")
	templates, err := s.repo.GetActiveTemplates(ctx)
	if err != nil {
		s.log.Error("Failed to get active templates", err)
		return nil, errors.DatabaseError("failed to retrieve active templates", err)
	}
	return templates, nil
}

func (s *posterTemplateSubService) UpdateTemplate(ctx context.Context, template *models.PosterTemplate) error {
	s.log.Info("Updating template", "id", template.ID)
	if err := s.repo.UpdateTemplate(ctx, template); err != nil {
		s.log.Error("Failed to update template", err, "id", template.ID)
		return errors.DatabaseError("failed to update template", err)
	}
	return nil
}

func (s *posterTemplateSubService) DeleteTemplate(ctx context.Context, id uint) error {
	s.log.Info("Deleting template", "id", id)
	if err := s.repo.DeleteTemplate(ctx, id); err != nil {
		s.log.Error("Failed to delete template", err, "id", id)
		return errors.DatabaseError("failed to delete template", err)
	}
	return nil
}
