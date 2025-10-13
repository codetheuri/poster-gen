package services

import (
	"context"
	"encoding/json"
	stdErrors "errors" // For errors.Is
	"fmt"
	"strings"
	"time"

	dto "github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto" // Import DTOs
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
	GeneratePoster(ctx context.Context, userID uint, orderID uint, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error)
	GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error)
	UpdatePoster(ctx context.Context, id uint, input *dto.PosterInput) error
	DeletePoster(ctx context.Context, id uint) error
}

type posterSubService struct {
	repo         repositories.PosterSubRepository
	validator    *validators.Validator
	log          logger.Logger
	templateRepo repositories.PosterTemplateSubRepository
}

// NewPosterSubService constructor
func NewPosterSubService(repo repositories.PosterSubRepository, validator *validators.Validator, log logger.Logger, templateRepo repositories.PosterTemplateSubRepository) PosterSubService {
	return &posterSubService{
		repo:         repo,
		validator:    validator,
		log:          log,
		templateRepo: templateRepo,
	}
}

// func (s *posterSubService) GeneratePoster(ctx context.Context, userID uint, orderID uint, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error) {
//     s.log.Info("Generating poster", "user_id", userID, "order_id", orderID, "template_id", templateID)

//     // Validate input
//     validationErrors := s.validator.Struct(input)
//     if validationErrors != nil {
//         s.log.Warn("Validation failed for poster input", validationErrors)
//         return nil, errors.ValidationError("invalid poster input", nil, validationErrors)
//     }

//     // Basic payment validation
//     if input.PaymentType == "mpesa" && input.TillNumber != "" && len(input.TillNumber) != 6 {
//         return nil, errors.ValidationError("till number must be 6 digits for M-Pesa", nil, nil)
//     }

//     // Generate QR code
//     qrData := fmt.Sprintf("https://yourapp.com/pay?till=%s&order=%d", input.TillNumber, orderID)
//     qrFile := fmt.Sprintf("posters/qr_%d.png", time.Now().Unix())
//     if err := qrcode.WriteFile(qrData, qrcode.Medium, 256, qrFile); err != nil {
//         s.log.Error("Failed to generate QR code", err)
//         return nil, errors.InternalServerError("failed to generate QR code", err)
//     }

//     // Generate PDF
//     pdf := gofpdf.New("P", "mm", "A4", "")
//     pdf.AddPage()
//     pdf.SetFont("Arial", "B", 16)
//     pdf.Cell(40, 10, fmt.Sprintf("Pay to: %s", input.BusinessName))
//     pdf.Ln(10)
//     pdf.Cell(40, 10, fmt.Sprintf("%s Till: %s", input.PaymentType, input.TillNumber))
//     pdf.Image(qrFile, 10, 30, 50, 50, false, "", 0, "")
//     pdfOutput := fmt.Sprintf("posters/%s_%d.pdf", input.BusinessName, time.Now().Unix())
//     if err := pdf.OutputFileAndClose(pdfOutput); err != nil {
//         s.log.Error("Failed to generate PDF", err)
//         return nil, errors.InternalServerError("failed to generate PDF", err)
//     }

//     // Create poster record
//     poster := &models.Poster{
//         UserID:        userID,
//         OrderID:       orderID,
//         TemplateID:    templateID,
//         BusinessName:  input.BusinessName,
//         PhoneNumber:   input.PhoneNumber,
//         PaymentType:   input.PaymentType,
//         TillNumber:    &input.TillNumber,
//         PaybillNumber: &input.PaybillNumber,
//         AccountNumber: &input.AccountNumber,
//         AccountName:   &input.AccountName,
//         BankName:      &input.BankName,
//         LogoURL:       input.LogoURL,
//         QRType:        "payment",
//         QRData:        qrData,
//         PDFURL:        pdfOutput,
//         Status:        "completed",
//     }

//     if err := s.repo.CreatePoster(ctx, poster); err != nil {
//         s.log.Error("Failed to save poster to database", err)
//         return nil, errors.DatabaseError("failed to save poster", err)
//     }

//	    // Return DTO response
//	    resp := &dto.PosterResponse{
//	        ID:          poster.ID,
//	        UserID:      poster.UserID,
//	        OrderID:     poster.OrderID,
//	        TemplateID:  poster.TemplateID,
//	        BusinessName: poster.BusinessName,
//	        PDFURL:      poster.PDFURL,
//	        Status:      poster.Status,
//	    }
//	    return resp, nil
//	}
func (s *posterSubService) GeneratePoster(ctx context.Context, userID uint, orderID uint, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error) {
	s.log.Info("Generating poster", "user_id", userID, "order_id", orderID, "template_id", templateID)

	// Validate input
	validationErrors := s.validator.Struct(input)
	if validationErrors != nil {
		s.log.Warn("Validation failed for poster input", validationErrors)
		return nil, errors.ValidationError("invalid poster input", nil, validationErrors)
	}

	// Fetch template
	template, err := s.templateRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Template not found", "template_id", templateID)
			return nil, errors.NotFoundError("template not found", err)
		}
		s.log.Error("Failed to get template", err, "template_id", templateID)
		return nil, errors.DatabaseError("failed to retrieve template", err)
	}

	// Parse layout
	var layout struct {
		Background string           `json:"background"`
		Elements   []map[string]any `json:"elements"`
	}
	if err := json.Unmarshal([]byte(template.Layout), &layout); err != nil {
		s.log.Error("Failed to parse template layout", err)
		return nil, errors.InternalServerError("invalid template layout", err)
	}

	// Generate QR code
	qrData := fmt.Sprintf("https://yourapp.com/pay?till=%s&order=%d", input.TillNumber, orderID)
	qrFile := fmt.Sprintf("posters/qr_%d.png", time.Now().Unix())
	if err := qrcode.WriteFile(qrData, qrcode.Medium, 256, qrFile); err != nil {
		s.log.Error("Failed to generate QR code", err)
		return nil, errors.InternalServerError("failed to generate QR code", err)
	}

	// Generate PDF with template
	pdf := gofpdf.New("P", "mm", template.Type, "") // Use template type (e.g., "A4")
	pdf.AddPage()
	if layout.Background != "" {
		pdf.Image(layout.Background, 0, 0, 210, 297, false, "", 0, "") // A4 size in mm
	}

	// Replace placeholders and add elements
	data := map[string]string{
		"business_name": input.BusinessName,
		"payment_type":  input.PaymentType,
		"till_number":   input.TillNumber,
		"qr_code":       qrFile,
	}
	for _, elem := range layout.Elements {
		switch elem["type"] {
		case "text":
			content, _ := elem["content"].(string)
			for key, val := range data {
				content = strings.ReplaceAll(content, "{"+key+"}", val)
			}
			pdf.SetFont(elem["font"].(string), "B", elem["size"].(float64))
			pdf.Text(elem["x"].(float64), elem["y"].(float64), content)
		case "image":
			pdf.Image(elem["content"].(string), elem["x"].(float64), elem["y"].(float64), elem["width"].(float64), elem["height"].(float64), false, "", 0, "")
		}
	}

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
		Status:        "completed",
	}

	if err := s.repo.CreatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to save poster to database", err)
		return nil, errors.DatabaseError("failed to save poster", err)
	}

	resp := &dto.PosterResponse{
		ID:           poster.ID,
		UserID:       poster.UserID,
		OrderID:      poster.OrderID,
		TemplateID:   poster.TemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}
	return resp, nil
}

func (s *posterSubService) GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error) {
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
	resp := &dto.PosterResponse{
		ID:           poster.ID,
		UserID:       poster.UserID,
		OrderID:      poster.OrderID,
		TemplateID:   poster.TemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}
	return resp, nil
}

func (s *posterSubService) UpdatePoster(ctx context.Context, id uint, input *dto.PosterInput) error {
	s.log.Info("Updating poster", "id", id)

	// Fetch existing poster
	poster, err := s.repo.GetPosterByID(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("Poster not found", "id", id)
			return errors.NotFoundError("poster not found", err)
		}
		s.log.Error("Failed to get poster for update", err, "id", id)
		return errors.DatabaseError("failed to retrieve poster", err)
	}

	// Update fields
	poster.BusinessName = input.BusinessName // Extend with other fields as needed
	if err := s.repo.UpdatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to update poster", err, "id", id)
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
