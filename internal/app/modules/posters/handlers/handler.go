package handlers
import (
   //	"context"
	//"net/http"
	postersServices "github.com/codetheuri/poster-gen/internal/app/modules/posters/services"
	"github.com/codetheuri/poster-gen/pkg/logger"
	//"github.com/codetheuri/poster-gen/pkg/web"
	//"github.com/codetheuri/poster-gen/internal/app/modules/posters/models"
)

type PostersHandler struct {
       postersService postersServices.PostersService
	   log logger.Logger
}

// constructor for PostersHandler
func NewPostersHandler(postersService postersServices.PostersService, log logger.Logger) *PostersHandler {
	return &PostersHandler{
		postersService: postersService,
		log: log,
	}
}

//example handler method

// func (h *PostersHandler) GetPostersByID(w http.ResponseWriter, r *http.Request) {
// 	h.log.Info("GetPostersByID handler invoked")
// 	// For example, to decode a request body into a model from this module:
// 	// var item models.Posters
// 	// if err := json.NewDecoder(r.Body).Decode(&item); err != nil { /* handle error */ }
// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Hello from Posters handler!"})
// }
