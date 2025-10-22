package services

import (
	"context"
)

// LogoSubService defines the interface for our new service.
type LogoSubService interface {
	GetLogos(ctx context.Context) ([]Logo, error)
}

type logoSubService struct {
	// This service is so simple it doesn't need any dependencies right now.
}

// NewLogoSubService is the constructor.
func NewLogoSubService() LogoSubService {
	return &logoSubService{}
}

// GetLogos simply returns the hardcoded data from our library.
func (s *logoSubService) GetLogos(ctx context.Context) ([]Logo, error) {
	logos := GetLogos()
	return logos, nil
}
