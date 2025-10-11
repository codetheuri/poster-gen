package repositories
 import (
    //  "context"
	//  "github.com/codetheuri/poster-gen/internal/app/modules/posters/models"
	  "github.com/codetheuri/poster-gen/pkg/logger"
	  "gorm.io/gorm"
 )
//define interface
     type PostersRepository interface {
	// Example method:
	// CreatePosters(ctx context.Context, posters *models.Posters) error
}
	type postersRepository struct {
	db *gorm.DB
	log logger.Logger 
	}
	// repo constructor
func NewPostersRepository(db *gorm.DB, log logger.Logger) PostersRepository {	
	return &postersRepository{
		db: db,
		log: log,
	}
	}
	//repos methods
	//eg.
	// func (r *gormPostersRepository) CreatePosters(ctx context.Context, posters *models.Posters) error {
// 	r.log.Info("CreatePosters repository ")
// 	// Placeholder for actual database logic
// 	return nil
// }

