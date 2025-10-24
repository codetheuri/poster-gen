package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"errors"
	// "math" // No longer needed directly if converting uint64 carefully

	postersDTO "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	// Import the main service package
	postersServices "github.com/codetheuri/poster-gen/internal/app/posters/services"
	appErrors "github.com/codetheuri/poster-gen/pkg/errors"

	// Removed tokenPkg import as getUserIDFromContext is removed/commented
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/codetheuri/poster-gen/pkg/web"
	"github.com/go-chi/chi"
)

// PostersHandler interface includes all methods handled by this package.
type PostersHandler interface {
	GeneratePoster(w http.ResponseWriter, r *http.Request)
	GetPosterByID(w http.ResponseWriter, r *http.Request)
	// UpdatePoster(w http.ResponseWriter, r *http.Request) // Placeholder
	// DeletePoster(w http.ResponseWriter, r *http.Request) // Placeholder
	GetActiveTemplates(w http.ResponseWriter, r *http.Request)
	CreateTemplate(w http.ResponseWriter, r *http.Request)
	GetTemplateByID(w http.ResponseWriter, r *http.Request)
	UpdateTemplate(w http.ResponseWriter, r *http.Request)
	DeleteTemplate(w http.ResponseWriter, r *http.Request)
	GetLogos(w http.ResponseWriter, r *http.Request) 
	CreateLayout(w http.ResponseWriter, r *http.Request)
	ListLayouts(w http.ResponseWriter, r *http.Request)
	CreateAsset(w http.ResponseWriter, r *http.Request)
	ListAssets(w http.ResponseWriter, r *http.Request)

}

type postersHandler struct {
	// Use the main service aggregator struct
	service   *postersServices.PosterService
	log       logger.Logger
	validator *validators.Validator
}

// NewPostersHandler constructor accepts the main service aggregator.
func NewPostersHandler(service *postersServices.PosterService, log logger.Logger, validator *validators.Validator) PostersHandler {
	return &postersHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

// GeneratePoster handles requests to create a new poster anonymously.
func (h *postersHandler) GeneratePoster(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GeneratePoster request")

	var input postersDTO.PosterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode GeneratePoster request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for GeneratePoster request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}
	templateIDStr := r.URL.Query().Get("template_id")

	templateID, err := strconv.ParseUint(templateIDStr, 10, 32) // Use 32 for uint
	if err != nil {
		h.log.Warn("Handler: Invalid or missing template_id", err, "template_id", templateIDStr)
		web.RespondError(w, appErrors.ValidationError("invalid template_id", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Call the correct sub-service via the main service aggregator
	poster, err := h.service.PosterSvc.GeneratePoster(ctx, uint(templateID), &input)
	if err != nil {
		h.log.Error("Handler: Failed to generate poster through service", err)
		h.handleAppError(w, err, "generate poster")
		return
	}

	h.log.Info("Handler: Poster generated successfully", "poster_id", poster.ID)
	// --- Send response using the correct payload structure ---
	// Your web.RespondData likely creates {"datapayload": {"data": poster}, "alertify": ...}
	web.RespondData(w, http.StatusCreated, poster, "Poster generated successfully", web.WithSuccessType("toast"))
}

// GetLogos handles requests for the predefined logo library.
func (h *postersHandler) GetLogos(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetLogos request")

	ctx := r.Context()
	// Call the LogoSvc sub-service
	logos, err := h.service.LogoSvc.GetLogos(ctx)
	if err != nil {
		// Logo service currently has no error paths, but handle just in case
		h.log.Error("Handler: Failed to get logos from service", err)
		h.handleAppError(w, err, "get logos")
		return
	}

	h.log.Info("Handler: Logos retrieved successfully", "count", len(logos))
	// --- Send response using the correct payload structure ---
	// Your web.RespondListData likely creates {"listdatapayload": {"data": logos, "pagination": null}}
	web.RespondListData(w, http.StatusOK, logos, nil) // Respond with the list
}

// GetPosterByID retrieves details of a specific generated poster.
func (h *postersHandler) GetPosterByID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetPosterByID request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Use 32 for uint
	if err != nil {
		h.log.Warn("Handler: Invalid poster ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid poster ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Call the correct sub-service
	poster, err := h.service.PosterSvc.GetPosterByID(ctx, uint(id))
	if err != nil {
		h.log.Error("Handler: Failed to get poster by ID", err, "id", id)
		h.handleAppError(w, err, "get poster")
		return
	}

	h.log.Info("Handler: Poster retrieved successfully", "poster_id", poster.ID)
	web.RespondData(w, http.StatusOK, poster, "Poster retrieved successfully", web.WithoutSuccess())
}

// GetActiveTemplates retrieves all currently active template profiles.
func (h *postersHandler) GetActiveTemplates(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetActiveTemplates request")

	ctx := r.Context()
	// Call the correct sub-service
	templates, err := h.service.PosterTemplateSvc.GetActiveTemplates(ctx)
	if err != nil {
		h.log.Error("Handler: Failed to get active templates", err)
		h.handleAppError(w, err, "get active templates")
		return
	}

	h.log.Info("Handler: Active templates retrieved successfully", "count", len(templates))
	// Use RespondListData for arrays
	web.RespondListData(w, http.StatusOK, templates, nil)
}

// CreateTemplate handles creation of new template profiles (admin).
func (h *postersHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received CreateTemplate request")

	var input postersDTO.TemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode CreateTemplate request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	if validationErrors := h.validator.Struct(input); validationErrors != nil {
		h.log.Warn("Handler: Validation failed for CreateTemplate request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Call the correct sub-service
	template, err := h.service.PosterTemplateSvc.CreateTemplate(ctx, &input)
	if err != nil {
		h.log.Error("Handler: Failed to create template through service", err)
		h.handleAppError(w, err, "create template")
		return
	}

	h.log.Info("Handler: Template created successfully", "template_id", template.ID)
	web.RespondData(w, http.StatusCreated, template, "Template created successfully", web.WithSuccessType("toast"))
}

// GetTemplateByID retrieves details of a specific template profile (admin).
func (h *postersHandler) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetTemplateByID request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Use 32 for uint
	if err != nil {
		h.log.Warn("Handler: Invalid template ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid template ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Call the correct sub-service
	template, err := h.service.PosterTemplateSvc.GetTemplateByID(ctx, uint(id))
	if err != nil {
		h.log.Error("Handler: Failed to get template by ID", err, "id", id)
		h.handleAppError(w, err, "get template")
		return
	}

	h.log.Info("Handler: Template retrieved successfully", "template_id", template.ID)
	web.RespondData(w, http.StatusOK, template, "Template retrieved successfully", web.WithoutSuccess())
}

// UpdateTemplate handles updates to existing template profiles (admin).
func (h *postersHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received UpdateTemplate request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Use 32 for uint
	if err != nil {
		h.log.Warn("Handler: Invalid template ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid template ID format", nil, nil), http.StatusBadRequest)
		return
	}

	var input postersDTO.TemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode UpdateTemplate request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	// Note: Validation might need adjustment for updates (e.g., allow partial updates)
	if validationErrors := h.validator.Struct(input); validationErrors != nil {
		h.log.Warn("Handler: Validation failed for UpdateTemplate request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Call the correct sub-service
	if err := h.service.PosterTemplateSvc.UpdateTemplate(ctx, uint(id), &input); err != nil {
		h.log.Error("Handler: Failed to update template", err, "id", id)
		h.handleAppError(w, err, "update template")
		return
	}

	h.log.Info("Handler: Template updated successfully", "template_id", id)
	web.RespondMessage(w, http.StatusOK, "Template updated successfully", "success", "toast")
}

// DeleteTemplate handles deletion of template profiles (admin).
func (h *postersHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received DeleteTemplate request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Use 32 for uint
	if err != nil {
		h.log.Warn("Handler: Invalid template ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid template ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Call the correct sub-service
	if err := h.service.PosterTemplateSvc.DeleteTemplate(ctx, uint(id)); err != nil {
		h.log.Error("Handler: Failed to delete template", err, "id", id)
		h.handleAppError(w, err, "delete template")
		return
	}

	h.log.Info("Handler: Template deleted successfully", "template_id", id)
	web.RespondMessage(w, http.StatusNoContent, "Template deleted successfully", "success", "toast")
}

func (h *postersHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received CreateLayout request")
	var input postersDTO.LayoutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}
	if validationErrors := h.validator.Struct(input); validationErrors != nil {
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Assume LayoutSvc has a CreateLayout method
	layout, err := h.service.LayoutSvc.CreateLayout(ctx, &input) // You need to implement CreateLayout in LayoutService
	if err != nil {
		h.log.Error("Handler: Failed to create layout", err)
		h.handleAppError(w, err, "create layout")
		return
	}
	h.log.Info("Handler: Layout created successfully", "layout_id", layout.ID)
	web.RespondData(w, http.StatusCreated, layout, "Layout created successfully", web.WithSuccessType("toast"))
}

func (h *postersHandler) ListLayouts(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received ListLayouts request")
	ctx := r.Context()
	// Assume LayoutSvc has a ListLayouts method
	layouts, err := h.service.LayoutSvc.ListLayouts(ctx) // You need to implement ListLayouts in LayoutService
	if err != nil {
		h.log.Error("Handler: Failed to list layouts", err)
		h.handleAppError(w, err, "list layouts")
		return
	}
	h.log.Info("Handler: Layouts listed successfully", "count", len(layouts))
	web.RespondListData(w, http.StatusOK, layouts, nil)
}

// --- Implementations for Assets ---

func (h *postersHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received CreateAsset request")
	var input postersDTO.AssetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}
	if validationErrors := h.validator.Struct(input); validationErrors != nil {
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Assume AssetSvc has a CreateAsset method
	asset, err := h.service.AssetSvc.CreateAsset(ctx, &input) // You need to implement CreateAsset in AssetService
	if err != nil {
		h.log.Error("Handler: Failed to create asset", err)
		h.handleAppError(w, err, "create asset")
		return
	}
	h.log.Info("Handler: Asset created successfully", "asset_id", asset.ID)
	web.RespondData(w, http.StatusCreated, asset, "Asset created successfully", web.WithSuccessType("toast"))
}

func (h *postersHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received ListAssets request")
	assetType := r.URL.Query().Get("type") // Optional: filter by type (e.g., /assets?type=logo)

	ctx := r.Context()
	// Assume AssetSvc has a ListAssets method
	assets, err := h.service.AssetSvc.ListAssets(ctx, assetType) // You need to implement ListAssets in AssetService
	if err != nil {
		h.log.Error("Handler: Failed to list assets", err)
		h.handleAppError(w, err, "list assets")
		return
	}
	h.log.Info("Handler: Assets listed successfully", "count", len(assets))
	web.RespondListData(w, http.StatusOK, assets, nil)
}
// handleAppError centralizes error response logic.
func (h *postersHandler) handleAppError(w http.ResponseWriter, err error, action string) {
	var appErr appErrors.AppError
	if errors.As(err, &appErr) {
		h.log.Error(fmt.Sprintf("Handler: Application error during %s", action), err, "code", appErr.Code())
		switch appErr.Code() {
		case "AUTH_ERROR":
			web.RespondError(w, appErr, http.StatusUnauthorized)
		case "NOT_FOUND":
			web.RespondError(w, appErr, http.StatusNotFound)
		case "VALIDATION_ERROR":
			web.RespondError(w, appErr, http.StatusBadRequest)
		case "DATABASE_ERROR":
			web.RespondError(w, appErr, http.StatusInternalServerError)
		case "INTERNAL_SERVER_ERROR":
			web.RespondError(w, appErr, http.StatusInternalServerError)
		case "CONFLICT_ERROR":
			web.RespondError(w, appErr, http.StatusConflict)
		case "PAYMENT_ERROR":
			web.RespondError(w, appErr, http.StatusPaymentRequired)
		case "BAD_REQUEST":
			web.RespondError(w, appErr, http.StatusBadRequest)
		default:
			web.RespondError(w, appErrors.InternalServerError(
				fmt.Sprintf("an unexpected application error occurred with code %s", appErr.Code()), appErr), http.StatusInternalServerError)
		}
	} else {
		h.log.Error(fmt.Sprintf("Handler: Unknown error during %s", action), err)
		web.RespondError(w, appErrors.InternalServerError("an unknown error occurred", err), http.StatusInternalServerError)
	}
}


