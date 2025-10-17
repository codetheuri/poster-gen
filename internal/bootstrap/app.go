package bootstrap

import (
	"fmt"
	"net"
	"net/http"

	"github.com/codetheuri/poster-gen/config"
	modules "github.com/codetheuri/poster-gen/internal/app"

	authModule "github.com/codetheuri/poster-gen/internal/app/auth"
	postersModule "github.com/codetheuri/poster-gen/internal/app/posters"
	router "github.com/codetheuri/poster-gen/internal/app/routers"

	"github.com/codetheuri/poster-gen/internal/platform/database"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/middleware"
	"github.com/codetheuri/poster-gen/pkg/validators"
	// "github.com/codetheuri/poster-gen/pkg/validators"
)

// initiliazes and start the application
func Run(cfg *config.Config, log logger.Logger) error {
	//db
	db, err := database.NewGoRMDB(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	//initialize the router

	//initilialize app components
	appValidator := validators.NewValidator()

	//application modules
	var appModules []modules.Module
	authMod := authModule.NewModule(db, log, appValidator, cfg)
	// Example of adding a new module))
	appModules = append(appModules, authModule.NewModule(db, log, appValidator, cfg))                     // Example of adding a new module
	appModules = append(appModules, postersModule.NewModule(db, log, appValidator, authMod.TokenService)) // Example of adding a new module

	//register routes from all modules
	mainRouter := router.NewRouter(log)
	for _, module := range appModules {
		module.RegisterRoutes(mainRouter)
	}

	//middleware
	var handler http.Handler = mainRouter
	handler = middleware.Logger(log)(handler)
	handler = middleware.Recovery(log)(handler)
	handler = middleware.RequestID()(handler)

	//Start Server
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	// get the actual address assigned (useful if port was 0)
	actualAddr := ln.Addr().(*net.TCPAddr)
	log.Info(fmt.Sprintf("Server is listening on port %d", actualAddr.Port))

	if err := http.Serve(ln, handler); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil

}
