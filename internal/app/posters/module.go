package posters

import (
	postersHandlers "github.com/codetheuri/poster-gen/internal/app/posters/handlers"
	postersRepositories "github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	postersServices "github.com/codetheuri/poster-gen/internal/app/posters/services"
	tokenPkg "github.com/codetheuri/poster-gen/pkg/auth/token"
	"github.com/codetheuri/poster-gen/internal/app/routers"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/middleware"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/gorm"
)

// Module represents the Posters module.
type Module struct {
	Handler      postersHandlers.PostersHandler
	log          logger.Logger
	TokenService tokenPkg.TokenService
}

// NewModule initializes the Posters module.
func NewModule(db *gorm.DB, log logger.Logger, validator *validators.Validator, tokenService tokenPkg.TokenService) *Module {
	repos := postersRepositories.NewPosterRepository(db, log)
	services := postersServices.NewPosterService(repos, validator, repos.PosterTemplateRepo, log)
	handler := postersHandlers.NewPostersHandler(services, log, validator)

	return &Module{
		Handler:      handler,
		log:          log,
		TokenService: tokenService,
	}
}

// RegisterRoutes registers the routes for the Posters module.
func (m *Module) RegisterRoutes(r router.Router) {
	m.log.Info("Registering Posters module routes...")

	// Public routes
	r.Group(func(r router.Router) {
		r.Get("/posters/templates", m.Handler.GetActiveTemplates)
		r.Post("/posters/generate", m.Handler.GeneratePoster)
		r.Get("/posters/{id}", m.Handler.GetPosterByID)
		r.Put("/posters/{id}", m.Handler.UpdatePoster)
		r.Delete("/posters/{id}", m.Handler.DeletePoster)
		r.Get("/logos", m.Handler.GetLogos)
	})

	// Authenticated routes
	r.Group(func(r router.Router) {
		r.Use(middleware.Authenticator(m.TokenService, m.log))

		r.Post("/posters/orders", m.Handler.CreateOrder)
		r.Post("/posters/orders/{id}/pay", m.Handler.ProcessPayment)
		r.Get("/posters/orders/{id}", m.Handler.GetOrderByID)
		r.Put("/posters/orders/{id}", m.Handler.UpdateOrder)
		r.Delete("/posters/orders/{id}", m.Handler.DeleteOrder)
		r.Post("/posters/templates", m.Handler.CreateTemplate)
		r.Get("/posters/templates/{id}", m.Handler.GetTemplateByID)
		r.Patch("/posters/templates/{id}", m.Handler.UpdateTemplate)
		r.Delete("/posters/templates/{id}", m.Handler.DeleteTemplate)
	})

	m.log.Info("Posters module routes registered.")
}
