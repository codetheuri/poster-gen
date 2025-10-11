package services
import (
	//"context"
	postersRepositories "github.com/codetheuri/poster-gen/internal/app/modules/posters/repositories"
	//"github.com/codetheuri/poster-gen/internal/app/modules/posters/models"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
)
	//define methods
type PostersService interface {
	// Define mthods here : eg
	//GetPostersByID(ctx context.Context, id int) (*models.Posters, error)
}

type postersService struct {
	Repo postersRepositories.PostersRepository
	validator *validators.Validator
	log logger.Logger
}

//service constructor
func NewPostersService(Repo postersRepositories.PostersRepository, validator *validators.Validator, log logger.Logger) PostersService {
	return &postersService{
		Repo: Repo,
		validator: validator,
		log: log,
	}
}

//methods
// func (s *postersService) GetPostersByID(ctx context.Context, id uint) (*models.Posters, error) {
// 	s.log.Info("GetPostersByID service invoked")
// 	// Placeholder for actual logic
// 	return nil, nil
// }
