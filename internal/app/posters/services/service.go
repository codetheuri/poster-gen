package services

import (
	"context" // Needed for service method signatures
	// Needed for error formatting
	// Need DTOs for input parameters
	dto "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	// Need Models for creating/returning data
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	posterRepositories "github.com/codetheuri/poster-gen/internal/app/posters/repositories"

	// Need errors package
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	// Needed for gorm.ErrRecordNotFound
)

// PosterService is the main service aggregator for the posters module.
type PosterService struct {
	PosterTemplateSvc PosterTemplateSubService
	PosterSvc         PosterSubService
	LogoSvc           LogoSubService
	LayoutSvc         LayoutSubService
	AssetSvc          AssetSubService
}

// NewPosterService constructor for the main service aggregator.
func NewPosterService(
	repos *posterRepositories.PosterRepository,
	validator *validators.Validator,
	log logger.Logger,
) *PosterService {
	templatesDir := "./templates"
	outputDir := "./posters"

	return &PosterService{
		PosterTemplateSvc: NewPosterTemplateSubService(repos.PosterTemplateRepo, repos.LayoutRepo, validator, log),
		PosterSvc:         NewPosterSubService(repos.PosterRepo, repos.PosterTemplateRepo, repos.LayoutRepo, repos.AssetRepo, validator, log, templatesDir, outputDir),
		LogoSvc:           NewLogoSubService(),
		LayoutSvc:         NewLayoutSubService(repos.LayoutRepo, log),
		AssetSvc:          NewAssetSubService(repos.AssetRepo, log),
		// OrderSvc:          NewOrderSubService(repos.OrderRepo, validator, log), // Keep commented if needed
	}
}

type LayoutSubService interface {
	CreateLayout(ctx context.Context, input *dto.LayoutInput) (*models.Layout, error)
	ListLayouts(ctx context.Context) ([]*models.Layout, error)
	// Add GetLayoutByID etc. if needed later
}
type layoutSubService struct {
	repo repositories.LayoutRepository
	log  logger.Logger
}

func NewLayoutSubService(repo repositories.LayoutRepository, log logger.Logger) LayoutSubService {
	return &layoutSubService{repo: repo, log: log}
}

// CreateLayout handles the business logic for creating a layout.
func (s *layoutSubService) CreateLayout(ctx context.Context, input *dto.LayoutInput) (*models.Layout, error) {
	s.log.Info("Creating layout", "name", input.Name)
	// Basic validation (can add more complex checks if needed)
	if input.Name == "" || input.FilePath == "" {
		return nil, errors.ValidationError("layout name and file path are required", nil, nil)
	}

	layout := &models.Layout{
		Name:     input.Name,
		FilePath: input.FilePath,
	}

	// Assuming the LayoutRepository has a CreateLayout method
	err := s.repo.CreateLayout(ctx, layout) // You need to add CreateLayout to the LayoutRepository interface and implementation
	if err != nil {
		s.log.Error("Failed to create layout in repository", err)
		// Check for specific DB errors like unique constraints if needed
		return nil, errors.DatabaseError("failed to save layout", err)
	}
	s.log.Info("Layout created successfully", "id", layout.ID)
	return layout, nil
}

// ListLayouts handles retrieving all layouts.
func (s *layoutSubService) ListLayouts(ctx context.Context) ([]*models.Layout, error) {
	s.log.Info("Listing layouts")
	// Assuming the LayoutRepository has a ListLayouts method
	layouts, err := s.repo.ListLayouts(ctx) // You need to add ListLayouts to the LayoutRepository interface and implementation
	if err != nil {
		s.log.Error("Failed to list layouts from repository", err)
		return nil, errors.DatabaseError("failed to retrieve layouts", err)
	}
	s.log.Info("Layouts listed successfully", "count", len(layouts))
	return layouts, nil
}

// --- Asset Service Implementation ---

// AssetSubService interface defines methods for asset operations.
type AssetSubService interface {
	CreateAsset(ctx context.Context, input *dto.AssetInput) (*models.Asset, error)
	ListAssets(ctx context.Context, assetType string) ([]*models.Asset, error)
	// Add GetAssetByID etc. if needed later
}
type assetSubService struct {
	repo repositories.AssetRepository
	log  logger.Logger
}

func NewAssetSubService(repo repositories.AssetRepository, log logger.Logger) AssetSubService {
	return &assetSubService{repo: repo, log: log}
}

// CreateAsset handles the business logic for creating an asset.
func (s *assetSubService) CreateAsset(ctx context.Context, input *dto.AssetInput) (*models.Asset, error) {
	s.log.Info("Creating asset", "name", input.Name, "type", input.Type)
	// Basic validation
	if input.Name == "" || input.Type == "" || input.Data == "" {
		return nil, errors.ValidationError("asset name, type, and data are required", nil, nil)
	}

	asset := &models.Asset{
		Name:         input.Name,
		Type:         input.Type,
		Data:         input.Data,
		DefaultColor: input.DefaultColor,
	}

	// Assuming the AssetRepository has a CreateAsset method
	err := s.repo.CreateAsset(ctx, asset) // You need to add CreateAsset to the AssetRepository interface and implementation
	if err != nil {
		s.log.Error("Failed to create asset in repository", err)
		return nil, errors.DatabaseError("failed to save asset", err)
	}
	s.log.Info("Asset created successfully", "id", asset.ID)
	return asset, nil
}

// ListAssets retrieves assets, optionally filtering by type.
func (s *assetSubService) ListAssets(ctx context.Context, assetType string) ([]*models.Asset, error) {
	s.log.Info("Listing assets", "type_filter", assetType)
	var assets []*models.Asset
	var err error

	// Assuming AssetRepository has appropriate methods
	if assetType != "" {
		assets, err = s.repo.GetAssetsByType(ctx, assetType) // You already have GetAssetsByType
	} else {
		assets, err = s.repo.ListAllAssets(ctx) // You need to add ListAllAssets to the AssetRepository interface and implementation
	}

	if err != nil {
		s.log.Error("Failed to list assets from repository", err, "type_filter", assetType)
		return nil, errors.DatabaseError("failed to retrieve assets", err)
	}
	s.log.Info("Assets listed successfully", "count", len(assets))
	return assets, nil
}
