package posters

import (
	postersHandlers "github.com/codetheuri/poster-gen/internal/app/modules/posters/handlers"
	postersRepositories "github.com/codetheuri/poster-gen/internal/app/modules/posters/repositories"
	postersServices "github.com/codetheuri/poster-gen/internal/app/modules/posters/services"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
)

// Module represents the Posters module.
type Module struct {
	Handlers *postersHandlers.PostersHandler

}

// NewModule initializes  Posters module.
func NewModule(db *gorm.DB, log logger.Logger, validator *validators.Validator) *Module {
     repo := postersRepositories.NewPostersRepository(db, log)
	 service := postersServices.NewPostersService(repo, validator, log)
	 handler := postersHandlers.NewPostersHandler(service, log)

	return &Module{
		Handlers: handler,	
}
}

// RegisterRoutes registers the routes for the Posters module.
func (m *Module) RegisterRoutes(r chi.Router) {
	// Register the routes for the posters module
	r.Route("/posters", func(r chi.Router) {
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Module posters is working!"))
	})
		//r.Post("/", m.Handlers.CreatePosters)
		//r.Get("/", m.Handlers.GetAllPosterss)
		
	})
}
