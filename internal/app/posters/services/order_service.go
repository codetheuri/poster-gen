package services

import (
	"context"
	stdErrors "errors"
	"fmt"
	"time"

	dto "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"

	// "github.com/codetheuri/poster-gen/pkg/mpesa"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/gorm"
)

type OrderSubService interface {
	CreateOrder(ctx context.Context, userID uint, input *dto.OrderInput) (*dto.OrderResponse, error)
	ProcessPayment(ctx context.Context, orderID uint, phoneNumber string) error
	GetOrderByID(ctx context.Context, id uint) (*dto.OrderResponse, error)
	UpdateOrder(ctx context.Context, id uint, input *dto.OrderInput) error
	DeleteOrder(ctx context.Context, id uint) error
}

type orderSubService struct {
	repo      repositories.OrderSubRepository
	validator *validators.Validator
	log       logger.Logger
	// mpesa     mpesa.Client
}

func NewOrderSubService(repo repositories.OrderSubRepository, validator *validators.Validator, log logger.Logger) OrderSubService {
	// mpesaClient := mpesa.NewClient("consumer_key", "consumer_secret", "sandbox") // Example
	return &orderSubService{
		repo:      repo,
		validator: validator,
		log:       log,
		// mpesa:     mpesaClient,
	}
}

func (s *orderSubService) CreateOrder(ctx context.Context, userID uint, input *dto.OrderInput) (*dto.OrderResponse, error) {
	s.log.Info("Creating order", "user_id", userID, "total_amount", input.TotalAmount)

	validationErrors := s.validator.Struct(input)
	if validationErrors != nil {
		s.log.Warn("Validation failed for order input", validationErrors)
		return nil, errors.ValidationError("invalid order input", nil, validationErrors)
	}

	orderNumber := fmt.Sprintf("ORD-%d-%d", userID, time.Now().Unix())
	order := &models.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		TotalAmount: input.TotalAmount,
		Status:      "pending",
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.log.Error("Failed to save order to database", err)
		return nil, errors.DatabaseError("failed to save order", err)
	}

	resp := &dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}
	return resp, nil
}

func (s *orderSubService) ProcessPayment(ctx context.Context, orderID uint, phoneNumber string) error {
	s.log.Info("Processing payment", "order_id", orderID, "phone_number", phoneNumber)

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Order not found", "order_id", orderID)
			return errors.NotFoundError("order not found", err)
		}
		s.log.Error("Failed to get order for payment", err, "order_id", orderID)
		return errors.DatabaseError("failed to retrieve order", err)
	}

	if order.Status != "pending" {
		return errors.BadRequestError("order is not in pending state", nil)
	}

	// receipt, err := s.mpesa.STKPush(phoneNumber, order.TotalAmount, order.OrderNumber)
	// if err != nil {
	// 	s.log.Error("Failed to process M-Pesa payment", err)
	// 	return errors.PaymentError("failed to process payment", err)
	// }

	// order.Status = "paid"
	// order.MpesaReceipt = &receipt
	// if err := s.repo.UpdateOrder(ctx, order); err != nil {
	// 	s.log.Error("Failed to update order after payment", err, "order_id", orderID)
	// 	return errors.DatabaseError("failed to update order", err)
	// }

	return nil
}

func (s *orderSubService) GetOrderByID(ctx context.Context, id uint) (*dto.OrderResponse, error) {
	s.log.Info("Getting order by ID", "id", id)
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Order not found", "id", id)
			return nil, errors.NotFoundError("order not found", err)
		}
		s.log.Error("Failed to get order", err, "id", id)
		return nil, errors.DatabaseError("failed to retrieve order", err)
	}
	resp := &dto.OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		OrderNumber:  order.OrderNumber,
		TotalAmount:  order.TotalAmount,
		Status:       order.Status,
		MpesaReceipt: &order.MpesaReceipt,
	}
	return resp, nil
}

func (s *orderSubService) UpdateOrder(ctx context.Context, id uint, input *dto.OrderInput) error {
	s.log.Info("Updating order", "id", id)

	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Order not found", "id", id)
			return errors.NotFoundError("order not found", err)
		}
		s.log.Error("Failed to get order for update", err, "id", id)
		return errors.DatabaseError("failed to retrieve order", err)
	}

	order.TotalAmount = input.TotalAmount // Extend with other fields as needed
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		s.log.Error("Failed to update order", err, "id", id)
		return errors.DatabaseError("failed to update order", err)
	}
	return nil
}

func (s *orderSubService) DeleteOrder(ctx context.Context, id uint) error {
	s.log.Info("Deleting order", "id", id)
	if err := s.repo.DeleteOrder(ctx, id); err != nil {
		s.log.Error("Failed to delete order", err, "id", id)
		return errors.DatabaseError("failed to delete order", err)
	}
	return nil
}
