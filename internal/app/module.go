package app

import "github.com/codetheuri/poster-gen/internal/app/routers"

// Module defines the contract that all application modules must follow.
type Module interface {
	// RegisterRoutes now requires our generic, framework-agnostic router.
	RegisterRoutes(r router.Router)
}
