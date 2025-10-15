package services

import (
	posterRepositories "github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
)

// PosterService is the main service aggregator for the posters module
type PosterService struct {
	PosterService         PosterSubService
	PosterTemplateService PosterTemplateSubService
	OrderService          OrderSubService

}

// NewPosterService constructor for the main service
func NewPosterService(
	repos *posterRepositories.PosterRepository,
	validator *validators.Validator,
	templateRepo posterRepositories.PosterTemplateSubRepository,
	log logger.Logger,
) *PosterService {
	return &PosterService{
		PosterService:         NewPosterSubService(repos.PosterRepo, validator, log, templateRepo, "./templates", "./posters"),
		PosterTemplateService: NewPosterTemplateSubService(repos.PosterTemplateRepo, validator, log),
		OrderService:          NewOrderSubService(repos.OrderRepo, validator, log),
	}
}
