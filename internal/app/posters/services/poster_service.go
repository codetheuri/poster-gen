package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/codetheuri/poster-gen/internal/app/posters/handlers/dto"
	"github.com/codetheuri/poster-gen/internal/app/posters/models"
	"github.com/codetheuri/poster-gen/internal/app/posters/repositories"
	"github.com/codetheuri/poster-gen/pkg/errors"
	"github.com/codetheuri/poster-gen/pkg/logger"
	"github.com/codetheuri/poster-gen/pkg/validators"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type PosterSubService interface {
	GeneratePoster(ctx context.Context, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error)
	GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error)
	UpdatePoster(ctx context.Context, id uint, input *dto.PosterInput) error
	DeletePoster(ctx context.Context, id uint) error
}

type posterSubService struct {
	repo         repositories.PosterSubRepository
	validator    *validators.Validator
	log          logger.Logger
	templateRepo repositories.PosterTemplateSubRepository
	templatesDir string
	outputDir    string
}

// NewPosterSubService constructor
func NewPosterSubService(
	repo repositories.PosterSubRepository,
	validator *validators.Validator,
	log logger.Logger,
	templateRepo repositories.PosterTemplateSubRepository,
	templatesDir string,
	outputDir string,
) PosterSubService {

	// Create directories if they don't exist
	os.MkdirAll(templatesDir, 0755)
	os.MkdirAll(outputDir, 0755)

	return &posterSubService{
		repo:         repo,
		validator:    validator,
		log:          log,
		templateRepo: templateRepo,
		templatesDir: templatesDir,
		outputDir:    outputDir,
	}
}

// GeneratePoster - Main function that orchestrates the entire process
func (s *posterSubService) GeneratePoster(ctx context.Context, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error) {
	s.log.Info("Generating poster", "template_id", templateID)

	// Validate input
	validationErrors := s.validator.Struct(input)
	if validationErrors != nil {
		s.log.Warn("Validation failed for poster input", validationErrors)
		return nil, errors.ValidationError("invalid poster input", nil, validationErrors)
	}

	// Fetch template
	templateRecord, err := s.templateRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Warn("Template not found", "template_id", templateID)
			return nil, errors.NotFoundError("template not found", err)
		}
		s.log.Error("Failed to get template", err, "template_id", templateID)
		return nil, errors.DatabaseError("failed to retrieve template", err)
	}

	// Prepare template data for HTML rendering
	templateData := map[string]interface{}{
		"business_name": input.BusinessName,
	}
	for key, value := range input.Data {
		templateData[key] = value

		// Smartly add split versions for number boxing
		if key == "agent_number" || key == "store_number" || key == "till_number" {
			if strValue, ok := value.(string); ok {
				templateData[key+"Split"] = strings.Split(strValue, "")
			}
		}
	}

	htmlContent, err := s.renderHTMLTemplate(templateData, templateRecord.Layout)
	if err != nil {
		s.log.Error("Failed to render HTML template", err)
		return nil, errors.InternalServerError("failed to render template", err)
	}
	// Generate PDF using HTML template
	pdfPath, err := s.htmlToPDF(ctx, htmlContent, input.BusinessName)
	if err != nil {
		s.log.Error("Failed to generate PDF from HTML", err)
		return nil, errors.InternalServerError("failed to generate PDF", err)
	}
	dynamicDataJSON, err := json.Marshal(input.Data)
	if err != nil {
		s.log.Error("Failed to marshal dynamic data to JSON", err)
		return nil, errors.InternalServerError("failed to marshal dynamic data", err)
	}
	// Create poster record
	poster := &models.Poster{
		UserID:       nil,
		OrderID:      nil,
		TemplateID:   templateID,
		BusinessName: input.BusinessName,
		DynamicData:  dynamicDataJSON, // <-- Store the entire dynamic data map as JSON
		PDFURL:       pdfPath,
		Status:       "completed",
	}

	if err := s.repo.CreatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to save poster to database", err)
		return nil, errors.DatabaseError("failed to save poster", err)
	}

	resp := &dto.PosterResponse{
		ID:           poster.ID,
		TemplateID:   poster.TemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}
	return resp, nil
}

// renderHTMLTemplate - Renders HTML template with data
func (s *posterSubService) renderHTMLTemplate(data map[string]interface{}, templateName string) (string, error) {
	templatePath := filepath.Join(s.templatesDir, templateName)

	// Read template file
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}

	// Parse and execute template
	tmpl, err := template.New("poster").Parse(string(templateBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// htmlToPDF - Converts HTML content to PDF using wkhtmltopdf

func (s *posterSubService) htmlToPDF(ctx context.Context, htmlContent string, businessName string) (string, error) {
	// Create a new context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfBuffer []byte
	pdfPath := filepath.Join(s.outputDir, fmt.Sprintf("%s_%d.pdf", businessName, time.Now().Unix()))

	err := chromedp.Run(ctx,
		// 1. Navigate to a blank page with the rendered HTML content
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		// 2. Wait until the page is loaded
		chromedp.WaitVisible("body", chromedp.ByQuery),
		// chromedp.Sleep(2*time.Second), // Give some time for all resources to load
		chromedp.Evaluate(`new Promise(resolve => {
        if (document.readyState === 'complete') {
            resolve();
        } else {
            window.addEventListener('load', resolve);
        }
    })`, nil),
	    chromedp.Sleep(4 * time.Second), // Extra wait to ensure all resources are loaded)
		chromedp.ActionFunc(func(ctx context.Context) error {

			buf, _, err := page.PrintToPDF().
				WithLandscape(true).
				WithPreferCSSPageSize(true).
				WithPaperWidth(11.7).
				WithPaperHeight(8.28).
				WithPrintBackground(true).
				WithMarginTop(0).
				WithMarginBottom(0).
			
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuffer = buf
			return nil
		}),
	)

	if err != nil {
		return "", fmt.Errorf("chromedp failed: %w", err)
	}

	// Write the PDF buffer to a file
	if err := os.WriteFile(pdfPath, pdfBuffer, 0644); err != nil {
		return "", fmt.Errorf("failed to write PDF file: %w", err)
	}

	return pdfPath, nil
}

// generateQRCode - Generates QR code image file
func (s *posterSubService) generateQRCode(data string) (string, error) {
	qrFile := filepath.Join(s.outputDir, fmt.Sprintf("qr_%d.png", time.Now().UnixNano()))

	if err := qrcode.WriteFile(data, qrcode.Medium, 256, qrFile); err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrFile, nil
}

// Other interface methods...
func (s *posterSubService) GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error) {
	poster, err := s.repo.GetPosterByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PosterResponse{
		ID:           poster.ID,
		TemplateID:   poster.TemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}, nil
}

func (s *posterSubService) UpdatePoster(ctx context.Context, id uint, input *dto.PosterInput) error {
	// Implementation for updating poster
	return nil
}

func (s *posterSubService) DeletePoster(ctx context.Context, id uint) error {
	return s.repo.DeletePoster(ctx, id)
}
