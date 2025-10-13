package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"errors"
	"math"

	postersDTO "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	postersServices "github.com/codetheuri/poster-gen/internal/app/posters/services"
	appErrors "github.com/codetheuri/poster-gen/pkg/errors"
	tokenPkg "github.com/codetheuri/poster-gen/pkg/auth/token"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/codetheuri/poster-gen/pkg/web"
	"github.com/go-chi/chi"
)

type PostersHandler interface {
	GeneratePoster(w http.ResponseWriter, r *http.Request)
	GetPosterByID(w http.ResponseWriter, r *http.Request)
	UpdatePoster(w http.ResponseWriter, r *http.Request)
	DeletePoster(w http.ResponseWriter, r *http.Request)
	GetActiveTemplates(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
	ProcessPayment(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	UpdateOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
	CreateTemplate(w http.ResponseWriter, r *http.Request)  // New
	GetTemplateByID(w http.ResponseWriter, r *http.Request) // New
	UpdateTemplate(w http.ResponseWriter, r *http.Request)  // New
	DeleteTemplate(w http.ResponseWriter, r *http.Request)  // New
}

type postersHandler struct {
	service   *postersServices.PosterService
	log       logger.Logger
	validator *validators.Validator
}

// NewPostersHandler constructor
func NewPostersHandler(service *postersServices.PosterService, log logger.Logger, validator *validators.Validator) PostersHandler {
	return &postersHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

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

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		h.log.Warn("Handler: UserID not found in context for GeneratePoster")
		web.RespondError(w, appErrors.AuthError("authentication context missing", nil), http.StatusUnauthorized)
		return
	}

	orderIDStr := r.URL.Query().Get("order_id")
	templateIDStr := r.URL.Query().Get("template_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID > math.MaxUint {
		h.log.Warn("Handler: Invalid or missing order_id", err, "order_id", orderIDStr)
		web.RespondError(w, appErrors.ValidationError("invalid order_id", nil, nil), http.StatusBadRequest)
		return
	}
	templateID, err := strconv.ParseUint(templateIDStr, 10, 64)
	if err != nil || templateID > math.MaxUint {
		h.log.Warn("Handler: Invalid or missing template_id", err, "template_id", templateIDStr)
		web.RespondError(w, appErrors.ValidationError("invalid template_id", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	poster, err := h.service.PosterService.GeneratePoster(ctx, uint(userID), uint(orderID), uint(templateID), &input)
	if err != nil {
		h.log.Error("Handler: Failed to generate poster through service", err)
		h.handleAppError(w, err, "generate poster")
		return
	}

	h.log.Info("Handler: Poster generated successfully", "poster_id", poster.ID)
	web.RespondData(w, http.StatusCreated, poster, "Poster generated successfully", web.WithSuccessType("toast"))
}

func (h *postersHandler) GetPosterByID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetPosterByID request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid poster ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid poster ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	poster, err := h.service.PosterService.GetPosterByID(ctx, uint(id))
	if err != nil {
		h.log.Error("Handler: Failed to get poster by ID", err, "id", id)
		h.handleAppError(w, err, "get poster")
		return
	}

	h.log.Info("Handler: Poster retrieved successfully", "poster_id", poster.ID)
	web.RespondData(w, http.StatusOK, poster, "Poster retrieved successfully", web.WithoutSuccess())
}

func (h *postersHandler) UpdatePoster(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received UpdatePoster request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid poster ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid poster ID format", nil, nil), http.StatusBadRequest)
		return
	}

	var input postersDTO.PosterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode UpdatePoster request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for UpdatePoster request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.PosterService.UpdatePoster(ctx, uint(id), &input); err != nil {
		h.log.Error("Handler: Failed to update poster", err, "id", id)
		h.handleAppError(w, err, "update poster")
		return
	}

	h.log.Info("Handler: Poster updated successfully", "poster_id", id)
	web.RespondMessage(w, http.StatusOK, "Poster updated successfully", "success", "toast")
}

func (h *postersHandler) DeletePoster(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received DeletePoster request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid poster ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid poster ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.PosterService.DeletePoster(ctx, uint(id)); err != nil {
		h.log.Error("Handler: Failed to delete poster", err, "id", id)
		h.handleAppError(w, err, "delete poster")
		return
	}

	h.log.Info("Handler: Poster deleted successfully", "poster_id", id)
	web.RespondMessage(w, http.StatusNoContent, "Poster deleted successfully", "success", "toast")
}

func (h *postersHandler) GetActiveTemplates(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetActiveTemplates request")

	ctx := r.Context()
	templates, err := h.service.PosterTemplateService.GetActiveTemplates(ctx)
	if err != nil {
		h.log.Error("Handler: Failed to get active templates", err)
		h.handleAppError(w, err, "get active templates")
		return
	}

	h.log.Info("Handler: Active templates retrieved successfully", "count", len(templates))
	web.RespondListData(w, http.StatusOK, templates, nil)
}

func (h *postersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received CreateOrder request")

	var input postersDTO.OrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode CreateOrder request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for CreateOrder request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		h.log.Warn("Handler: UserID not found in context for CreateOrder")
		web.RespondError(w, appErrors.AuthError("authentication context missing", nil), http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	order, err := h.service.OrderService.CreateOrder(ctx, uint(userID), &input)
	if err != nil {
		h.log.Error("Handler: Failed to create order", err)
		h.handleAppError(w, err, "create order")
		return
	}

	h.log.Info("Handler: Order created successfully", "order_id", order.ID)
	web.RespondData(w, http.StatusCreated, order, "Order created successfully", web.WithSuccessType("toast"))
}

func (h *postersHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received ProcessPayment request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid order ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid order ID format", nil, nil), http.StatusBadRequest)
		return
	}

	var input struct {
		PhoneNumber string `json:"phone_number" validate:"required,len=10"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode ProcessPayment request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for ProcessPayment request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.OrderService.ProcessPayment(ctx, uint(id), input.PhoneNumber); err != nil {
		h.log.Error("Handler: Failed to process payment", err, "order_id", id)
		h.handleAppError(w, err, "process payment")
		return
	}

	h.log.Info("Handler: Payment processed successfully", "order_id", id)
	web.RespondMessage(w, http.StatusOK, "Payment processed successfully", "success", "toast")
}

func (h *postersHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetOrderByID request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid order ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid order ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	order, err := h.service.OrderService.GetOrderByID(ctx, uint(id))
	if err != nil {
		h.log.Error("Handler: Failed to get order by ID", err, "id", id)
		h.handleAppError(w, err, "get order")
		return
	}

	h.log.Info("Handler: Order retrieved successfully", "order_id", order.ID)
	web.RespondData(w, http.StatusOK, order, "Order retrieved successfully", web.WithoutSuccess())
}

func (h *postersHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received UpdateOrder request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid order ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid order ID format", nil, nil), http.StatusBadRequest)
		return
	}

	var input postersDTO.OrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode UpdateOrder request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for UpdateOrder request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.OrderService.UpdateOrder(ctx, uint(id), &input); err != nil {
		h.log.Error("Handler: Failed to update order", err, "id", id)
		h.handleAppError(w, err, "update order")
		return
	}

	h.log.Info("Handler: Order updated successfully", "order_id", id)
	web.RespondMessage(w, http.StatusOK, "Order updated successfully", "success", "toast")
}

func (h *postersHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received DeleteOrder request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid order ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid order ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.OrderService.DeleteOrder(ctx, uint(id)); err != nil {
		h.log.Error("Handler: Failed to delete order", err, "id", id)
		h.handleAppError(w, err, "delete order")
		return
	}

	h.log.Info("Handler: Order deleted successfully", "order_id", id)
	web.RespondMessage(w, http.StatusNoContent, "Order deleted successfully", "success", "toast")
}
func (h *postersHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received CreateTemplate request")

	var input postersDTO.TemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("Handler: Failed to decode CreateTemplate request", err)
		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
		return
	}

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for CreateTemplate request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	template, err := h.service.PosterTemplateService.CreateTemplate(ctx, &input)
	if err != nil {
		h.log.Error("Handler: Failed to create template through service", err)
		h.handleAppError(w, err, "create template")
		return
	}

	h.log.Info("Handler: Template created successfully", "template_id", template.ID)
	web.RespondData(w, http.StatusCreated, template, "Template created successfully", web.WithSuccessType("toast"))
}

func (h *postersHandler) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received GetTemplateByID request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid template ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid template ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	template, err := h.service.PosterTemplateService.GetTemplateByID(ctx, uint(id))
	if err != nil {
		h.log.Error("Handler: Failed to get template by ID", err, "id", id)
		h.handleAppError(w, err, "get template")
		return
	}

	h.log.Info("Handler: Template retrieved successfully", "template_id", template.ID)
	web.RespondData(w, http.StatusOK, template, "Template retrieved successfully", web.WithoutSuccess())
}

func (h *postersHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received UpdateTemplate request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
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

	validationErrors := h.validator.Struct(input)
	if validationErrors != nil {
		h.log.Warn("Handler: Validation failed for UpdateTemplate request", validationErrors)
		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.PosterTemplateService.UpdateTemplate(ctx, uint(id), &input); err != nil {
		h.log.Error("Handler: Failed to update template", err, "id", id)
		h.handleAppError(w, err, "update template")
		return
	}

	h.log.Info("Handler: Template updated successfully", "template_id", id)
	web.RespondMessage(w, http.StatusOK, "Template updated successfully", "success", "toast")
}

func (h *postersHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handler: Received DeleteTemplate request")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id > math.MaxUint {
		h.log.Warn("Handler: Invalid template ID format", err, "id", idStr)
		web.RespondError(w, appErrors.ValidationError("invalid template ID format", nil, nil), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.PosterTemplateService.DeleteTemplate(ctx, uint(id)); err != nil {
		h.log.Error("Handler: Failed to delete template", err, "id", id)
		h.handleAppError(w, err, "delete template")
		return
	}

	h.log.Info("Handler: Template deleted successfully", "template_id", id)
	web.RespondMessage(w, http.StatusNoContent, "Template deleted successfully", "success", "toast")
}
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

func getUserIDFromContext(ctx context.Context) (uint, bool) {
	return tokenPkg.GetUserIDFromContext(ctx)
}
