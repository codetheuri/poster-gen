package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	postersDTO "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	postersServices "github.com/codetheuri/poster-gen/internal/app/posters/services"
	appErrors "github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"

	// "github.com/codetheuri/poster-gen/pkg/pagination"
	"errors"
	"math"

	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/codetheuri/poster-gen/pkg/web"
	"github.com/go-chi/chi"
)

type PostersHandler interface {
	// GeneratePoster(w http.ResponseWriter, r *http.Request)
	GetPosterByID(w http.ResponseWriter, r *http.Request)
	// UpdatePoster(w http.ResponseWriter, r *http.Request)
	DeletePoster(w http.ResponseWriter, r *http.Request)
	GetActiveTemplates(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
	ProcessPayment(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	// UpdateOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
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

// func (h *postersHandler) GeneratePoster(w http.ResponseWriter, r *http.Request) {
// 	h.log.Info("Handler: Received GeneratePoster request")

// 	var input postersDTO.PosterInput
// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.log.Warn("Handler: Failed to decode GeneratePoster request", err)
// 		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
// 		return
// 	}

// 	validationErrors := h.validator.Struct(input)
// 	if validationErrors != nil {
// 		h.log.Warn("Handler: Validation failed for GeneratePoster request", validationErrors)
// 		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
// 		return
// 	}

// 	userID, ok := getUserIDFromContext(r.Context())
// 	if !ok {
// 		h.log.Warn("Handler: UserID not found in context for GeneratePoster")
// 		web.RespondError(w, appErrors.AuthError("authentication context missing", nil), http.StatusUnauthorized)
// 		return
// 	}

// 	// Extract orderID and templateID from query params or path if needed
// 	orderIDStr := r.URL.Query().Get("order_id")
// 	templateIDStr := r.URL.Query().Get("template_id")
// 	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
// 	if err != nil || orderID > math.MaxUint {
// 		h.log.Warn("Handler: Invalid or missing order_id", err, "order_id", orderIDStr)
// 		web.RespondError(w, appErrors.ValidationError("invalid order_id", nil, nil), http.StatusBadRequest)
// 		return
// 	}
// 	templateID, err := strconv.ParseUint(templateIDStr, 10, 64)
// 	if err != nil || templateID > math.MaxUint {
// 		h.log.Warn("Handler: Invalid or missing template_id", err, "template_id", templateIDStr)
// 		web.RespondError(w, appErrors.ValidationError("invalid template_id", nil, nil), http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	poster, err := h.service.PosterService.GeneratePoster(ctx, uint(userID), uint(orderID), uint(templateID), &input)
// 	if err != nil {
// 		h.log.Error("Handler: Failed to generate poster through service", err)
// 		h.handleAppError(w, err, "generate poster")
// 		return
// 	}

// 	resp := postersDTO.PosterResponse{
// 		ID:           poster.ID,
// 		UserID:       poster.UserID,
// 		OrderID:      poster.OrderID,
// 		TemplateID:   poster.TemplateID,
// 		BusinessName: poster.BusinessName,
// 		PDFURL:       poster.PDFURL,
// 		Status:       poster.Status,
// 	}
// 	h.log.Info("Handler: Poster generated successfully", "poster_id", poster.ID)
// 	web.RespondData(w, http.StatusCreated, resp, "Poster generated successfully", web.WithSuccessType("toast"))
// }

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

	resp := postersDTO.PosterResponse{
		ID:           poster.ID,
		UserID:       poster.UserID,
		OrderID:      poster.OrderID,
		TemplateID:   poster.TemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}
	h.log.Info("Handler: Poster retrieved successfully", "poster_id", poster.ID)
	web.RespondData(w, http.StatusOK, resp, "Poster retrieved successfully", web.WithoutSuccess())
}

// func (h *postersHandler) UpdatePoster(w http.ResponseWriter, r *http.Request) {
// 	h.log.Info("Handler: Received UpdatePoster request")

// 	idStr := chi.URLParam(r, "id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil || id > math.MaxUint {
// 		h.log.Warn("Handler: Invalid poster ID format", err, "id", idStr)
// 		web.RespondError(w, appErrors.ValidationError("invalid poster ID format", nil, nil), http.StatusBadRequest)
// 		return
// 	}

// 	var input postersDTO.PosterInput
// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.log.Warn("Handler: Failed to decode UpdatePoster request", err)
// 		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
// 		return
// 	}

// 	validationErrors := h.validator.Struct(input)
// 	if validationErrors != nil {
// 		h.log.Warn("Handler: Validation failed for UpdatePoster request", validationErrors)
// 		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	poster := &models.Poster{ID: uint(id)}   // Partial update; adjust fields as needed
// 	poster.BusinessName = input.BusinessName // Example; extend with other fields
// 	if err := h.service.PosterService.UpdatePoster(ctx, poster); err != nil {
// 		h.log.Error("Handler: Failed to update poster", err, "id", id)
// 		h.handleAppError(w, err, "update poster")
// 		return
// 	}

// 	h.log.Info("Handler: Poster updated successfully", "poster_id", id)
// 	web.RespondMessage(w, http.StatusOK, "Poster updated successfully", "success", "toast")
// }

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

	resp := make([]postersDTO.TemplateResponse, len(templates))
	for i, t := range templates {
		resp[i] = postersDTO.TemplateResponse{
			ID:        t.ID,
			Name:      t.Name,
			Type:      t.Type,
			Price:     t.Price,
			Thumbnail: t.Thumbnail,
			IsActive:  t.IsActive,
		}
	}
	h.log.Info("Handler: Active templates retrieved successfully", "count", len(templates))
	web.RespondListData(w, http.StatusOK, resp, nil)
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
	order, err := h.service.OrderService.CreateOrder(ctx, uint(userID), input.TotalAmount)
	if err != nil {
		h.log.Error("Handler: Failed to create order", err)
		h.handleAppError(w, err, "create order")
		return
	}

	resp := postersDTO.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}
	h.log.Info("Handler: Order created successfully", "order_id", order.ID)
	web.RespondData(w, http.StatusCreated, resp, "Order created successfully", web.WithSuccessType("toast"))
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

	resp := postersDTO.OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		OrderNumber:  order.OrderNumber,
		TotalAmount:  order.TotalAmount,
		Status:       order.Status,
		MpesaReceipt: &order.MpesaReceipt,
	}
	h.log.Info("Handler: Order retrieved successfully", "order_id", order.ID)
	web.RespondData(w, http.StatusOK, resp, "Order retrieved successfully", web.WithoutSuccess())
}

// func (h *postersHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
// 	h.log.Info("Handler: Received UpdateOrder request")

// 	idStr := chi.URLParam(r, "id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil || id > math.MaxUint {
// 		h.log.Warn("Handler: Invalid order ID format", err, "id", idStr)
// 		web.RespondError(w, appErrors.ValidationError("invalid order ID format", nil, nil), http.StatusBadRequest)
// 		return
// 	}

// 	var input postersDTO.OrderInput
// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.log.Warn("Handler: Failed to decode UpdateOrder request", err)
// 		web.RespondError(w, appErrors.ValidationError("invalid request payload", err, nil), http.StatusBadRequest)
// 		return
// 	}

// 	validationErrors := h.validator.Struct(input)
// 	if validationErrors != nil {
// 		h.log.Warn("Handler: Validation failed for UpdateOrder request", validationErrors)
// 		web.RespondError(w, appErrors.ValidationError("validation failed", nil, validationErrors), http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	order := &models.Order{ID: uint(id), TotalAmount: input.TotalAmount} // Partial update
// 	if err := h.service.OrderService.UpdateOrder(ctx, order); err != nil {
// 		h.log.Error("Handler: Failed to update order", err, "id", id)
// 		h.handleAppError(w, err, "update order")
// 		return
// 	}

// 	h.log.Info("Handler: Order updated successfully", "order_id", id)
// 	web.RespondMessage(w, http.StatusOK, "Order updated successfully", "success", "toast")
// }

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
			web.RespondError(w, appErr, http.StatusPaymentRequired) // 402, adjust if needed
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

// getUserIDFromContext extracts the user ID from the request context (set by middleware)
func getUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("userID").(uint)
	return userID, ok
}
