package services

import (
	"time"

	authRepositories "github.com/codetheuri/poster-gen/internal/app/auth/repositories"

	//"github.com/codetheuri/poster-gen/internal/app/modules/auth/models"
	tokenPkg "github.com/codetheuri/poster-gen/pkg/auth/token"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
)

type AuthService struct {
	UserService  UserService
	TokenService tokenPkg.TokenService
}

// service constructor for all services
func NewAuthService(
	repos *authRepositories.AuthRepository,
	validator *validators.Validator,
	jwtSecret string,
	tokenTTL time.Duration,
	log logger.Logger) *AuthService {
	return &AuthService{
		UserService:  NewUserService(repos.UserRepo, validator, log),
		TokenService: NewJWTService(repos.RevokedTokenRepo, jwtSecret, tokenTTL, log),
	}
}
