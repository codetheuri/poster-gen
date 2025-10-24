package services

import (
	"context"
	"encoding/json"
	stdErrors "errors"

	dto "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	// Use the specific repository interface
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PosterTemplateSubService interface defines operations for managing template profiles.
type PosterTemplateSubService interface {
	CreateTemplate(ctx context.Context, input *dto.TemplateInput) (*dto.TemplateResponse, error)
	GetTemplateByID(ctx context.Context, id uint) (*dto.TemplateResponse, error)
	GetActiveTemplates(ctx context.Context) ([]*dto.TemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uint, input *dto.TemplateInput) error
	DeleteTemplate(ctx context.Context, id uint) error
}

type posterTemplateSubService struct {
	repo      repositories.PosterTemplateRepository // Uses the specific repo interface
	layoutRepo repositories.LayoutRepository       // Added Layout Repo dependency
	validator *validators.Validator
	log       logger.Logger
}

// NewPosterTemplateSubService constructor accepts necessary repositories.
func NewPosterTemplateSubService(repo repositories.PosterTemplateRepository, layoutRepo repositories.LayoutRepository, validator *validators.Validator, log logger.Logger) PosterTemplateSubService {
	return &posterTemplateSubService{
		repo:       repo,
		layoutRepo: layoutRepo, // Store layout repo
		validator:  validator,
		log:        log,
	}
}

// CreateTemplate creates a new poster template (customization profile).
func (s *posterTemplateSubService) CreateTemplate(ctx context.Context, input *dto.TemplateInput) (*dto.TemplateResponse, error) {
	s.log.Info("Creating poster template", "name", input.Name)

	if validationErrors := s.validator.Struct(input); validationErrors != nil {
		s.log.Warn("Validation failed for template input", validationErrors)
		return nil, errors.ValidationError("invalid template input", nil, validationErrors)
	}

	// Optional: Validate that the referenced LayoutID exists
	_, err := s.layoutRepo.GetLayoutByID(ctx, input.LayoutID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Warn("Referenced layout not found", "layout_id", input.LayoutID)
			// Return a validation error specific to the layout ID
			return nil, errors.ValidationError("invalid layout_id: layout not found", nil, map[string]string{"layout_id": "Referenced layout does not exist"})

		}
		s.log.Error("Failed to verify layout ID", err, "layout_id", input.LayoutID)
		return nil, errors.DatabaseError("failed to verify layout", err)
	}


	template := &models.PosterTemplate{
		Name:                 input.Name,
		Type:                 input.Type,
		LayoutID:             input.LayoutID, // Use LayoutID from input DTO
		Price:                input.Price,
		ThumbnailURL:         input.ThumbnailURL, // Corrected field name
		IsActive:             input.IsActive,
		RequiredFields:       datatypes.JSON(input.RequiredFields),
		DefaultCustomization: datatypes.JSON(input.DefaultCustomization), // Use correct field name
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		s.log.Error("Failed to save template to database", err)
		// Check for unique constraint violation on name
		// if strings.Contains(err.Error(), "UNIQUE constraint failed") { ... return ConflictError ... }
		return nil, errors.DatabaseError("failed to save template", err)
	}

	// Fetch again to ensure Layout info is populated for the response
	createdTemplate, err := s.repo.GetTemplateByID(ctx, template.ID)
	if err != nil {
		s.log.Error("Failed to fetch created template with layout", err, "template_id", template.ID)
		// Return basic response if refetch fails
		return &dto.TemplateResponse{ ID: template.ID /* Populate other known fields */ }, nil
	}


	return &dto.TemplateResponse{
		ID:                   createdTemplate.ID,
		Name:                 createdTemplate.Name,
		Type:                 createdTemplate.Type,
		LayoutID:             createdTemplate.LayoutID,
		LayoutFilePath:       createdTemplate.Layout.FilePath, // Get path from loaded Layout
		Price:                createdTemplate.Price,
		ThumbnailURL:         createdTemplate.ThumbnailURL,
		IsActive:             createdTemplate.IsActive,
		RequiredFields:       json.RawMessage(createdTemplate.RequiredFields),
		DefaultCustomization: json.RawMessage(createdTemplate.DefaultCustomization),
	}, nil
}

// GetTemplateByID retrieves a template including its layout file path.
func (s *posterTemplateSubService) GetTemplateByID(ctx context.Context, id uint) (*dto.TemplateResponse, error) {
	s.log.Info("Getting template by ID", "id", id)
	template, err := s.repo.GetTemplateByID(ctx, id) // Repo Preloads Layout
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Template not found", "id", id)
			return nil, errors.NotFoundError("template not found", err)
		}
		s.log.Error("Failed to get template", err, "id", id)
		return nil, errors.DatabaseError("failed to retrieve template", err)
	}

	layoutFilePath := ""
	if template.Layout.FilePath != "" {
		layoutFilePath = template.Layout.FilePath
	} else {
        s.log.Warn("Layout FilePath missing for fetched template", nil, "template_id", id, "layout_id", template.LayoutID)
        // Optionally try fetching Layout again if Preload might have failed silently
        // layout, layoutErr := s.layoutRepo.GetLayoutByID(ctx, template.LayoutID) ...
    }


	return &dto.TemplateResponse{
		ID:                   template.ID,
		Name:                 template.Name,
		Type:                 template.Type,
		LayoutID:             template.LayoutID,
		LayoutFilePath:       layoutFilePath,
		Price:                template.Price,
		ThumbnailURL:         template.ThumbnailURL,
		IsActive:             template.IsActive,
		RequiredFields:       json.RawMessage(template.RequiredFields),
		DefaultCustomization: json.RawMessage(template.DefaultCustomization),
	}, nil
}

// GetActiveTemplates retrieves all active templates including layout file paths.
func (s *posterTemplateSubService) GetActiveTemplates(ctx context.Context) ([]*dto.TemplateResponse, error) {
	s.log.Info("Getting active templates")
	templates, err := s.repo.GetActiveTemplates(ctx) // Repo Preloads Layouts
	if err != nil {
		s.log.Error("Failed to get active templates", err)
		return nil, errors.DatabaseError("failed to retrieve active templates", err)
	}
	resp := make([]*dto.TemplateResponse, len(templates))
	for i, t := range templates {
		layoutFilePath := ""
		if t.Layout.FilePath != "" {
			layoutFilePath = t.Layout.FilePath
		} else {
             s.log.Warn("Layout FilePath missing for active template in list", nil, "template_id", t.ID, "layout_id", t.LayoutID)
        }
		resp[i] = &dto.TemplateResponse{
			ID:                   t.ID,
			Name:                 t.Name,
			Type:                 t.Type,
			LayoutID:             t.LayoutID,
			LayoutFilePath:       layoutFilePath,
			Price:                t.Price,
			ThumbnailURL:         t.ThumbnailURL,
			IsActive:             t.IsActive,
			RequiredFields:       json.RawMessage(t.RequiredFields),
			DefaultCustomization: json.RawMessage(t.DefaultCustomization),
		}
	}
	return resp, nil
}

// UpdateTemplate updates an existing poster template.
func (s *posterTemplateSubService) UpdateTemplate(ctx context.Context, id uint, input *dto.TemplateInput) error {
	s.log.Info("Updating template", "id", id)

	// Validate input first
	if validationErrors := s.validator.Struct(input); validationErrors != nil {
		s.log.Warn("Validation failed for update template input", validationErrors)
		return errors.ValidationError("invalid update template input", nil, validationErrors)
	}

	template, err := s.repo.GetTemplateByID(ctx, id) // Fetch existing
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Template not found for update", "id", id)
			return errors.NotFoundError("template not found", err)
		}
		s.log.Error("Failed to get template for update", err, "id", id)
		return errors.DatabaseError("failed to retrieve template", err)
	}

	// Optional: Validate LayoutID exists if it's being changed
	if input.LayoutID != 0 && input.LayoutID != template.LayoutID {
        _, err := s.layoutRepo.GetLayoutByID(ctx, input.LayoutID)
        if err != nil {
            // Handle layout not found error similar to CreateTemplate
             return errors.ValidationError("invalid layout_id: layout not found", nil, map[string]string{"layout_id": "Referenced layout does not exist"})
        }
        template.LayoutID = input.LayoutID
	}


	// Apply updates selectively
	template.Name = input.Name // Assume required fields in DTO are always provided for update
	template.Type = input.Type
	template.Price = input.Price
	template.ThumbnailURL = input.ThumbnailURL
	template.IsActive = input.IsActive

	// Update JSON fields only if new data is provided in the input
	if len(input.RequiredFields) > 0 && string(input.RequiredFields) != "null" {
		template.RequiredFields = datatypes.JSON(input.RequiredFields)
	}
	if len(input.DefaultCustomization) > 0 && string(input.DefaultCustomization) != "null" {
		template.DefaultCustomization = datatypes.JSON(input.DefaultCustomization)
	}

	if err := s.repo.UpdateTemplate(ctx, template); err != nil {
		s.log.Error("Failed to update template in database", err, "id", id)
		return errors.DatabaseError("failed to update template", err)
	}
	s.log.Info("Template updated successfully", "id", id)
	return nil
}

// DeleteTemplate deletes a poster template.
func (s *posterTemplateSubService) DeleteTemplate(ctx context.Context, id uint) error {
	s.log.Info("Deleting template", "id", id)
	if err := s.repo.DeleteTemplate(ctx, id); err != nil {
		s.log.Error("Failed to delete template", err, "id", id)
		return errors.DatabaseError("failed to delete template", err)
	}
	s.log.Info("Template deleted successfully", "id", id)
	return nil
}

