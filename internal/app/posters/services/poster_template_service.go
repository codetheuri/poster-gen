package services

import (
	"context"
	stdErrors "errors"

	dto "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/gorm"
)

type PosterTemplateSubService interface {
	CreateTemplate(ctx context.Context, input *dto.TemplateInput) (*dto.TemplateResponse, error)
	GetTemplateByID(ctx context.Context, id uint) (*dto.TemplateResponse, error)
	GetActiveTemplates(ctx context.Context) ([]*dto.TemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uint, input *dto.TemplateInput) error
	DeleteTemplate(ctx context.Context, id uint) error
}

type posterTemplateSubService struct {
	repo      repositories.PosterTemplateSubRepository
	validator *validators.Validator
	log       logger.Logger
}

func NewPosterTemplateSubService(repo repositories.PosterTemplateSubRepository, validator *validators.Validator, log logger.Logger) PosterTemplateSubService {
	return &posterTemplateSubService{
		repo:      repo,
		validator: validator,
		log:       log,
	}
}

func (s *posterTemplateSubService) CreateTemplate(ctx context.Context, input *dto.TemplateInput) (*dto.TemplateResponse, error) {
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

	resp := &dto.TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      template.Type,
		Price:     template.Price,
		Thumbnail: template.Thumbnail,
		IsActive:  template.IsActive,
	}
	return resp, nil
}

func (s *posterTemplateSubService) GetTemplateByID(ctx context.Context, id uint) (*dto.TemplateResponse, error) {
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
	resp := &dto.TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      template.Type,
		Price:     template.Price,
		Thumbnail: template.Thumbnail,
		IsActive:  template.IsActive,
	}
	return resp, nil
}

func (s *posterTemplateSubService) GetActiveTemplates(ctx context.Context) ([]*dto.TemplateResponse, error) {
	s.log.Info("Getting active templates")
	templates, err := s.repo.GetActiveTemplates(ctx)
	if err != nil {
		s.log.Error("Failed to get active templates", err)
		return nil, errors.DatabaseError("failed to retrieve active templates", err)
	}
	resp := make([]*dto.TemplateResponse, len(templates))
	for i, t := range templates {
		resp[i] = &dto.TemplateResponse{
			ID:        t.ID,
			Name:      t.Name,
			Type:      t.Type,
			Price:     t.Price,
			Thumbnail: t.Thumbnail,
			IsActive:  t.IsActive,
		}
	}
	return resp, nil
}

func (s *posterTemplateSubService) UpdateTemplate(ctx context.Context, id uint, input *dto.TemplateInput) error {
	s.log.Info("Updating template", "id", id)

	template, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Template not found", "id", id)
			return errors.NotFoundError("template not found", err)
		}
		s.log.Error("Failed to get template for update", err, "id", id)
		return errors.DatabaseError("failed to retrieve template", err)
	}

	template.Name = input.Name
	template.Type = input.Type
	template.Price = input.Price
	template.Thumbnail = input.Thumbnail
	template.IsActive = input.IsActive

	if err := s.repo.UpdateTemplate(ctx, template); err != nil {
		s.log.Error("Failed to update template", err, "id", id)
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
