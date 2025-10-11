package services

import (
	"context"
	stdErrors "errors" // Import standard errors package
	"fmt"
	"time"

	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"
	// "github.com/codetheuri/poster-gen/pkg/mpesa" // Placeholder for M-Pesa SDK
	"github.com/codetheuri/poster-gen/pkg/validators"
	"gorm.io/gorm"
)

// OrderSubService interface for order operations
type OrderSubService interface {
	CreateOrder(ctx context.Context, userID uint, totalAmount int) (*models.Order, error)
	ProcessPayment(ctx context.Context, orderID uint, phoneNumber string) error
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, id uint) error
}

// OrderInput defines the input data for creating an order
type OrderInput struct {
	TotalAmount int `json:"total_amount" validate:"required,min=0"`
}

type orderSubService struct {
	repo      repositories.OrderSubRepository
	validator *validators.Validator
	log       logger.Logger
	// mpesa     mpesa.Client // Placeholder for M-Pesa client
}

// NewOrderSubService constructor
func NewOrderSubService(repo repositories.OrderSubRepository, validator *validators.Validator, log logger.Logger) OrderSubService {
	// Initialize M-Pesa client (e.g., with Daraja API credentials) - add in real impl
	// mpesaClient := mpesa.NewClient("consumer_key", "consumer_secret", "sandbox") // Example
	return &orderSubService{
		repo:      repo,
		validator: validator,
		log:       log,
		// mpesa:     mpesaClient,
	}
}

func (s *orderSubService) CreateOrder(ctx context.Context, userID uint, totalAmount int) (*models.Order, error) {
	s.log.Info("Creating order", "user_id", userID, "total_amount", totalAmount)

	input := &OrderInput{TotalAmount: totalAmount}
	validationErrors := s.validator.Struct(input)
	if validationErrors != nil {
		s.log.Warn("Validation failed for order input", validationErrors)
		return nil, errors.ValidationError("invalid order input", nil, validationErrors)
	}

	orderNumber := fmt.Sprintf("ORD-%d-%d", userID, time.Now().Unix())
	order := &models.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		TotalAmount: totalAmount,
		Status:      "pending",
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.log.Error("Failed to save order to database", err)
		return nil, errors.DatabaseError("failed to save order", err)
	}

	s.log.Info("Order created successfully", "order_id", order.ID)
	return order, nil
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

	// Placeholder for M-Pesa STK Push (replace with real SDK call)
	// receipt, err := s.mpesa.STKPush(phoneNumber, order.TotalAmount, order.OrderNumber)
	if err != nil {
		s.log.Error("Failed to process M-Pesa payment", err)
		return errors.PaymentError("failed to process payment", err)
	}

	order.Status = "paid"
	// order.MpesaReceipt = receipt
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		s.log.Error("Failed to update order after payment", err, "order_id", orderID)
		return errors.DatabaseError("failed to update order", err)
	}

	// s.log.Info("Payment processed successfully", "order_id", orderID, "receipt", receipt)
	return nil
}

func (s *orderSubService) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
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
	return order, nil
}

func (s *orderSubService) UpdateOrder(ctx context.Context, order *models.Order) error {
	s.log.Info("Updating order", "id", order.ID)
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		s.log.Error("Failed to update order", err, "id", order.ID)
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
