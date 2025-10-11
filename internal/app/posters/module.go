package posters

import (
	postersHandlers "github.com/codetheuri/poster-gen/internal/app/posters/handlers"
	postersRepositories "github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	postersServices "github.com/codetheuri/poster-gen/internal/app/posters/services"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

// Module represents the Posters module.
type Module struct {
	Handler postersHandlers.PostersHandler
	log     logger.Logger
}

// NewModule initializes the Posters module.
func NewModule(db *gorm.DB, log logger.Logger, validator *validators.Validator) *Module {
	repos := postersRepositories.NewPosterRepository(db, log)
	services := postersServices.NewPosterService(repos, validator, log)
	handler := postersHandlers.NewPostersHandler(services, log, validator)

	return &Module{
		Handler: handler,
		log:     log,
	}
}

// RegisterRoutes registers the routes for the Posters module.
func (m *Module) RegisterRoutes(r chi.Router) {
	m.log.Info("Registering Posters module routes...")

	// Public routes (e.g., listing templates)
	r.Group(func(r chi.Router) {
		r.Get("/posters/templates", m.Handler.GetActiveTemplates)
	})

	// Authenticated routes (require JWT)
	r.Group(func(r chi.Router) {
		// r.Use(middleware.Authenticator(nil, m.log)) // TODO: Inject TokenService from auth module
		// r.Post("/posters/generate", m.Handler.GeneratePoster)
		r.Get("/posters/{id}", m.Handler.GetPosterByID)
		// r.Put("/posters/{id}", m.Handler.UpdatePoster)
		r.Delete("/posters/{id}", m.Handler.DeletePoster)
		r.Post("/posters/orders", m.Handler.CreateOrder)
		r.Post("/posters/orders/{id}/pay", m.Handler.ProcessPayment)
		r.Get("/posters/orders/{id}", m.Handler.GetOrderByID)
		// r.Put("/posters/orders/{id}", m.Handler.UpdateOrder)
		r.Delete("/posters/orders/{id}", m.Handler.DeleteOrder)
	})

	m.log.Info("Posters module routes registered.")
}
