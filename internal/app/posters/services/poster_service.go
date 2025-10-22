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
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ... (Interface and struct definitions remain the same) ...
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

func NewPosterSubService(
	repo repositories.PosterSubRepository,
	validator *validators.Validator,
	log logger.Logger,
	templateRepo repositories.PosterTemplateSubRepository,
	templatesDir string,
	outputDir string,
) PosterSubService {
	// ... (constructor is fine) ...
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

// GeneratePoster - Reverted loop to standard form.
func (s *posterSubService) GeneratePoster(ctx context.Context, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error) {
	s.log.Info("Generating poster with dynamic template", "template_id", templateID)

	if validationErrors := s.validator.Struct(input); validationErrors != nil {
		return nil, errors.ValidationError("invalid poster input", nil, validationErrors)
	}
	templateRecord, err := s.templateRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("template not found", err)
		}
		return nil, errors.DatabaseError("failed to retrieve template", err)
	}
	templateData := make(map[string]interface{})
	var baseCustomizationData map[string]interface{}
	if len(templateRecord.CustomizationData) > 0 && string(templateRecord.CustomizationData) != "null" {
		if err := json.Unmarshal(templateRecord.CustomizationData, &baseCustomizationData); err != nil {
			var strData string
			if err2 := json.Unmarshal(templateRecord.CustomizationData, &strData); err2 == nil && strData != "" {

				if err3 := json.Unmarshal([]byte(strData), &baseCustomizationData); err3 != nil {
					s.log.Error("Failed to unmarshal nested customization data string", err3, "raw_data", strData)
					return nil, errors.InternalServerError("invalid template customization data format", err3)
				}
			} else {
				s.log.Error("Failed to unmarshal base customization data", err, "raw_data", string(templateRecord.CustomizationData))
				return nil, errors.InternalServerError("invalid template customization data", err)
			}
		}
	}

	// 2. Add base customization data to templateData.

	for key, value := range baseCustomizationData {
		if key == "header_logo_svg" {
			if svgStr, ok := value.(string); ok {
				templateData[key] = template.HTML(svgStr)
			} else {
				templateData[key] = value
			}
		} else {
			templateData[key] = value

		}
		// key is string here
	}

	// 3. MERGE: Add/Override with user's specific customization choices from input.
	if input.CustomizationData != nil {
		for k, v := range input.CustomizationData {
			if k == "header_logo_svg" {
				if svgStr, ok := v.(string); ok {
					templateData[k] = template.HTML(svgStr)
					
				}else{
					templateData[k] = v
				}
			} else {
				templateData[k] = v
			}
		}

	}

	// 4. Add user's dynamic form input data.

	if input.Data != nil {
		for key, value := range input.Data {
			templateData[key] = value
			if strings.HasSuffix(key, "_number") {
				if strValue, ok := value.(string); ok && strValue != "" {
					templateData[key+"Split"] = strings.Split(strValue, "")
				} else {
					templateData[key+"Split"] = []string{}
				}
			}
		}
	}

	// 5. Add business name.
	templateData["business_name"] = input.BusinessName

	htmlContent, err := s.renderHTMLTemplate(templateData, templateRecord.Layout)
	if err != nil {
		s.log.Error("Failed to render HTML template", err)
		return nil, errors.InternalServerError("failed to render template", err)
	}

	// pdfPath, err := s.htmlToPDF(ctx, htmlContent, input.BusinessName)
	pdfPath, err := s.renderToImage(ctx, htmlContent, input.BusinessName)
	if err != nil {
		s.log.Error("Failed to generate PDF from HTML", err)
		return nil, errors.InternalServerError("failed to generate PDF", err)
	}

	dynamicDataJSON, err := json.Marshal(input.Data)
	if err != nil {
		return nil, errors.InternalServerError("failed to marshal dynamic data", err)
	}

	poster := &models.Poster{
		UserID:       nil,
		OrderID:      nil,
		TemplateID:   templateID,
		BusinessName: input.BusinessName,
		DynamicData:  datatypes.JSON(dynamicDataJSON),
		PDFURL:       pdfPath,
		Status:       "completed",
	}

	if err := s.repo.CreatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to save poster to database", err)
		return nil, errors.DatabaseError("failed to save poster", err)
	}

	return &dto.PosterResponse{
		ID:           poster.ID,
		TemplateID:   poster.TemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}, nil
}

func (s *posterSubService) renderHTMLTemplate(data map[string]interface{}, templateName string) (string, error) {
	templatePath := filepath.Join(s.templatesDir, templateName)
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}
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

func (s *posterSubService) htmlToPDF(ctx context.Context, htmlContent string, businessName string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var pdfBuffer []byte
	pdfPath := filepath.Join(s.outputDir, fmt.Sprintf("%s_%d.pdf", businessName, time.Now().Unix()))
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.Evaluate(`new Promise(resolve => {
			if (document.readyState === 'complete') { resolve(); }
			else { window.addEventListener('load', resolve); }
		})`, nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
			    WithLandscape(true).
				WithPreferCSSPageSize(true).
				WithPrintBackground(true).
			WithPaperWidth(11.7).
				WithPaperHeight(8.27).
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
	if err := os.WriteFile(pdfPath, pdfBuffer, 0644); err != nil {
		return "", fmt.Errorf("failed to write PDF file: %w", err)
	}
	return pdfPath, nil
}

func (s *posterSubService) renderToImage(ctx context.Context, htmlContent string, businessName string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var imageBuffer []byte //
	imagePath := filepath.Join(s.outputDir, fmt.Sprintf("%s_%d.jpeg", businessName, time.Now().Unix()))

	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.Evaluate(`new Promise(resolve => {
			if (document.readyState === 'complete') { resolve(); }
			else { window.addEventListener('load', resolve); }
		})`, nil),

		// Capture the entire page as a screenshot (defaults to PNG)
		chromedp.FullScreenshot(&imageBuffer, 90), // 90 is image quality (for JPEG, not PNG)
	)
	if err != nil {
		return "", fmt.Errorf("chromedp failed to capture screenshot: %w", err)
	}
	if err := os.WriteFile(imagePath, imageBuffer, 0644); err != nil {
		return "", fmt.Errorf("failed to write image file: %w", err)
	}
	return imagePath, nil
}

func (s *posterSubService) GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error) {
	poster, err := s.repo.GetPosterByID(ctx, id)
	if err != nil {
		return nil, err // Assuming repo handles not found error appropriately
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
	// Implementation for updating poster - currently returns nil
	return nil
}
func (s *posterSubService) DeletePoster(ctx context.Context, id uint) error {
	return s.repo.DeletePoster(ctx, id)
}
