package auth

import (
	"github.com/codetheuri/poster-gen/config"
	authHandlers "github.com/codetheuri/poster-gen/internal/app/auth/handlers"
	authRepositories "github.com/codetheuri/poster-gen/internal/app/auth/repositories"
	authServices "github.com/codetheuri/poster-gen/internal/app/auth/services"
	tokenPkg "github.com/codetheuri/poster-gen/pkg/auth/token"
	"github.com/codetheuri/poster-gen/internal/app/routers"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/middleware"
	"github.com/codetheuri/poster-gen/pkg/validators"
	
	"gorm.io/gorm"
)

// Module represents the Auth module.
type Module struct {
	Handler      authHandlers.AuthHandler
	log          logger.Logger
	TokenService tokenPkg.TokenService
}

// NewModule initializes  Auth module.
func NewModule(db *gorm.DB, log logger.Logger, validator *validators.Validator, cfg *config.Config) *Module {
	repos := authRepositories.NewAuthRepository(db, log)

	jwtSecret := cfg.JWTSecret
	tokenTTL := cfg.AccessTokenTTL

	TokenService := authServices.NewJWTService(repos.RevokedTokenRepo, jwtSecret, tokenTTL, log)
	services := authServices.NewAuthService(repos, validator, jwtSecret, tokenTTL, log)
	handler := authHandlers.NewAuthHandler(services, log, validator)

	return &Module{
		Handler:      handler,
		TokenService: TokenService,
		log:          log,
	}
}

// RegisterRoutes registers the routes for the Auth module.
func (m *Module) RegisterRoutes(r router.Router) {
	m.log.Info("Registering Auth module routes...")

	r.Group(func(r router.Router) {
		r.Post("/auth/register", m.Handler.Register)
		r.Post("/auth/login", m.Handler.Login)
	})

	// Authenticated routes (will need middleware later)

	r.Group(func(r router.Router) {
		r.Use(middleware.Authenticator(m.TokenService, m.log))
		r.Get("/auth/profile/{id}", m.Handler.GetUserProfile)
		r.Put("/auth/users/{id}/change-password", m.Handler.ChangePassword)
		r.Delete("/auth/users/{id}", m.Handler.DeleteUser)
		r.Put("/auth/users/{id}/restore", m.Handler.RestoreUser)
		r.Post("/auth/logout", m.Handler.Logout)
		r.Get("/auth/users", m.Handler.GetUsers)
	})

	m.log.Info("Auth module routes registered.")
}
