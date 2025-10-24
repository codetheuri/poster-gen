package posters

import (
	postersHandlers "github.com/codetheuri/poster-gen/internal/app/posters/handlers"
	postersRepositories "github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	postersServices "github.com/codetheuri/poster-gen/internal/app/posters/services"
	"github.com/codetheuri/poster-gen/internal/app/routers"
	tokenPkg "github.com/codetheuri/poster-gen/pkg/auth/token"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/middleware"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/gorm"
)

type Module struct {
	Handler      postersHandlers.PostersHandler
	log          logger.Logger
	TokenService tokenPkg.TokenService // Keep if using authentication middleware
}

// NewModule initializes the Posters module using the aggregated service.
func NewModule(db *gorm.DB, log logger.Logger, validator *validators.Validator, tokenService tokenPkg.TokenService) *Module {
	// 1. Create the aggregated repository
	repos := postersRepositories.NewPosterRepository(db, log)
	// 2. Create the aggregated service, passing the aggregated repo
	services := postersServices.NewPosterService(repos, validator, log) // Pass only needed args
	// 3. Create the handler, passing the aggregated service
	handler := postersHandlers.NewPostersHandler(services, log, validator)

	return &Module{
		Handler:      handler,
		log:          log,
		TokenService: tokenService,
	}
}

// RegisterRoutes registers the routes for the Posters module using the generic router interface.
func (m *Module) RegisterRoutes(r router.Router) {
	m.log.Info("Registering Posters module routes...")

	r.Group(func(r router.Router) {
		r.Get("/posters/templates", m.Handler.GetActiveTemplates)
		r.Post("/posters/generate", m.Handler.GeneratePoster)
		r.Get("/posters/{id}", m.Handler.GetPosterByID) // Get generated poster details

		r.Get("/logos", m.Handler.GetLogos)
	})

	// Authenticated routes (Require JWT - For Admin/Management)
	r.Group(func(r router.Router) {
		// Apply authentication middleware for this group
		r.Use(middleware.Authenticator(m.TokenService, m.log))
		r.Post("/posters/templates", m.Handler.CreateTemplate)
		r.Get("/posters/templates/{id}", m.Handler.GetTemplateByID)
		r.Patch("/posters/templates/{id}", m.Handler.UpdateTemplate) // Use Patch for partial updates if applicable
		r.Delete("/posters/templates/{id}", m.Handler.DeleteTemplate)


		// Routes for managing Layouts (HTML structures) could go here
		r.Post("/layouts", m.Handler.CreateLayout)
		r.Get("/layouts", m.Handler.ListLayouts)
		
		r.Post("/assets", m.Handler.CreateAsset)
		r.Get("/assets", m.Handler.ListAssets)

	})

	m.log.Info("Posters module routes registered.")
}
