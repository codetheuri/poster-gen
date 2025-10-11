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
	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// PosterSubService interface for poster operations
type PosterSubService interface {
	GeneratePoster(ctx context.Context, userID uint, orderID uint, templateID uint, input *PosterInput) (*models.Poster, error)
	GetPosterByID(ctx context.Context, id uint) (*models.Poster, error)
	UpdatePoster(ctx context.Context, poster *models.Poster) error
	DeletePoster(ctx context.Context, id uint) error
}

// PosterInput defines the input data for generating a poster
type PosterInput struct {
	BusinessName  string `json:"business_name" validate:"required,max=100"`
	PhoneNumber   string `json:"phone_number" validate:"omitempty,max=15"`
	PaymentType   string `json:"payment_type" validate:"required,oneof=mpesa bank"`
	TillNumber    string `json:"till_number" validate:"omitempty,max=10"`
	PaybillNumber string `json:"paybill_number" validate:"omitempty,max=10"`
	AccountNumber string `json:"account_number" validate:"omitempty,max=20"`
	AccountName   string `json:"account_name" validate:"omitempty,max=100"`
	BankName      string `json:"bank_name" validate:"omitempty,max=50"`
	LogoURL       string `json:"logo_url" validate:"omitempty,url,max=255"`
}

type posterSubService struct {
	repo      repositories.PosterSubRepository
	validator *validators.Validator
	log       logger.Logger
}

// NewPosterSubService constructor
func NewPosterSubService(repo repositories.PosterSubRepository, validator *validators.Validator, log logger.Logger) PosterSubService {
	return &posterSubService{
		repo:      repo,
		validator: validator,
		log:       log,
	}
}

func (s *posterSubService) GeneratePoster(ctx context.Context, userID uint, orderID uint, templateID uint, input *PosterInput) (*models.Poster, error) {
	s.log.Info("Generating poster", "user_id", userID, "order_id", orderID, "template_id", templateID)

	// Validate input
	validationErrors := s.validator.Struct(input)
	if validationErrors != nil {
		s.log.Warn("Validation failed for poster input", validationErrors)
		return nil, errors.ValidationError("invalid poster input", nil, validationErrors)
	}

	// Basic payment validation (e.g., M-Pesa till number must be 6 digits if provided)
	if input.PaymentType == "mpesa" && input.TillNumber != "" && len(input.TillNumber) != 6 {
		return nil, errors.ValidationError("till number must be 6 digits for M-Pesa", nil, nil)
	}

	// Generate QR code (e.g., payment link)
	qrData := fmt.Sprintf("https://yourapp.com/pay?till=%s&order=%d", input.TillNumber, orderID)
	qrFile := fmt.Sprintf("posters/qr_%d.png", time.Now().Unix())
	if err := qrcode.WriteFile(qrData, qrcode.Medium, 256, qrFile); err != nil {
		s.log.Error("Failed to generate QR code", err)
		return nil, errors.InternalServerError("failed to generate QR code", err)
	}

	// Generate PDF (simplified)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Pay to: %s", input.BusinessName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("%s Till: %s", input.PaymentType, input.TillNumber))
	pdf.Image(qrFile, 10, 30, 50, 50, false, "", 0, "")
	pdfOutput := fmt.Sprintf("posters/%s_%d.pdf", input.BusinessName, time.Now().Unix())
	if err := pdf.OutputFileAndClose(pdfOutput); err != nil {
		s.log.Error("Failed to generate PDF", err)
		return nil, errors.InternalServerError("failed to generate PDF", err)
	}

	// Create poster record
	poster := &models.Poster{
		UserID:        userID,
		OrderID:       orderID,
		TemplateID:    templateID,
		BusinessName:  input.BusinessName,
		PhoneNumber:   input.PhoneNumber,
		PaymentType:   input.PaymentType,
		TillNumber:    &input.TillNumber,
		PaybillNumber: &input.PaybillNumber,
		AccountNumber: &input.AccountNumber,
		AccountName:   &input.AccountName,
		BankName:      &input.BankName,
		LogoURL:       input.LogoURL,
		QRType:        "payment",
		QRData:        qrData,
		PDFURL:        pdfOutput,
		ImageURL:      "", // Add image generation if needed
		Status:        "completed",
	}

	if err := s.repo.CreatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to save poster to database", err)
		return nil, errors.DatabaseError("failed to save poster", err)
	}

	s.log.Info("Poster generated successfully", "poster_id", poster.ID)
	return poster, nil
}

func (s *posterSubService) GetPosterByID(ctx context.Context, id uint) (*models.Poster, error) {
	s.log.Info("Getting poster by ID", "id", id)
	poster, err := s.repo.GetPosterByID(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Poster not found", "id", id)
			return nil, errors.NotFoundError("poster not found", err)
		}
		s.log.Error("Failed to get poster", err, "id", id)
		return nil, errors.DatabaseError("failed to retrieve poster", err)
	}
	return poster, nil
}

func (s *posterSubService) UpdatePoster(ctx context.Context, poster *models.Poster) error {
	s.log.Info("Updating poster", "id", poster.ID)
	if err := s.repo.UpdatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to update poster", err, "id", poster.ID)
		return errors.DatabaseError("failed to update poster", err)
	}
	return nil
}

func (s *posterSubService) DeletePoster(ctx context.Context, id uint) error {
	s.log.Info("Deleting poster", "id", id)
	if err := s.repo.DeletePoster(ctx, id); err != nil {
		s.log.Error("Failed to delete poster", err, "id", id)
		return errors.DatabaseError("failed to delete poster", err)
	}
	return nil
}
