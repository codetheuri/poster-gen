package repositories

import (
	//  "context"
	//  "github.com/codetheuri/poster-gen/internal/app/modules/auth/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"gorm.io/gorm"
)

type AuthRepository struct {
	UserRepo         UserRepository
	RevokedTokenRepo RevokedTokenRepository
}

// repo constructor
func NewAuthRepository(db *gorm.DB, log logger.Logger) *AuthRepository {
	return &AuthRepository{
		UserRepo:         NewUserRepository(db, log),
		RevokedTokenRepo: NewRevokedTokenRepository(db, log),
	}
}
