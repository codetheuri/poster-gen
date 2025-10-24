package repositories

import (
	"context"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)

// AssetRepository defines the interface for asset data operations.
type AssetRepository interface {
		CreateAsset(ctx context.Context, asset *models.Asset) error 
	ListAllAssets(ctx context.Context) ([]*models.Asset, error)
	GetAssetByID(ctx context.Context, id uint) (*models.Asset, error)
	GetAssetsByType(ctx context.Context, assetType string) ([]*models.Asset, error)

}

type assetRepository struct {
	db  *gorm.DB
	log logger.Logger
}

// NewAssetRepository creates a new AssetRepository.
func NewAssetRepository(db *gorm.DB, log logger.Logger) AssetRepository {
	return &assetRepository{db: db, log: log}
}
func (r *assetRepository) CreateAsset(ctx context.Context, asset *models.Asset) error {
	if err := r.db.WithContext(ctx).Create(asset).Error; err != nil {
		r.log.Error("Failed to create asset", err, "asset_name", asset.Name)
		return err 
	}
	return nil
}
func (r *assetRepository) ListAllAssets(ctx context.Context) ([]*models.Asset, error) {
	var assets []*models.Asset
	if err := r.db.WithContext(ctx).Find(&assets).Error; err != nil {
		r.log.Error("Failed to list all assets", err)
		return nil, err 
	}
	return assets, nil
}
func (r *assetRepository) GetAssetByID(ctx context.Context, id uint) (*models.Asset, error) {
	var asset models.Asset
	if err := r.db.WithContext(ctx).First(&asset, id).Error; err != nil {
		r.log.Error("Failed to get asset by ID", err, "asset_id", id)
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetAssetsByType(ctx context.Context, assetType string) ([]*models.Asset, error) {
	var assets []*models.Asset
	if err := r.db.WithContext(ctx).Where("type = ?", assetType).Find(&assets).Error; err != nil {
		r.log.Error("Failed to get assets by type", err, "asset_type", assetType)
		return nil, err
	}
	return assets, nil
}
