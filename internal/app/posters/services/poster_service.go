package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv" // Added for robust asset ID parsing
	"strings"
	"time"
	"unicode/utf8"

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
type RequiredFieldConfig struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	Type         string `json:"type"`
	Pattern      string `json:"pattern,omitempty"`
	MaxLength    int    `json:"maxLength,omitempty"`
	PatternTitle string `json:"patternTitle,omitempty"`
}
type PosterSubService interface {
	GeneratePoster(ctx context.Context, templateID uint, input *dto.PosterInput) (*dto.PosterResponse, error)
	GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error)
}

type posterSubService struct {
	repo         repositories.PosterSubRepository
	templateRepo repositories.PosterTemplateRepository
	layoutRepo   repositories.LayoutRepository
	assetRepo    repositories.AssetRepository
	validator    *validators.Validator
	log          logger.Logger
	templatesDir string
	outputDir    string
}

func NewPosterSubService(
	repo repositories.PosterSubRepository,
	templateRepo repositories.PosterTemplateRepository,
	layoutRepo repositories.LayoutRepository,
	assetRepo repositories.AssetRepository,
	validator *validators.Validator,
	log logger.Logger,
	templatesDir string,
	outputDir string,
) PosterSubService {
	os.MkdirAll(templatesDir, 0755)
	os.MkdirAll(outputDir, 0755)

	return &posterSubService{
		repo:         repo,
		templateRepo: templateRepo,
		layoutRepo:   layoutRepo,
		assetRepo:    assetRepo,
		validator:    validator,
		log:          log,
		templatesDir: templatesDir,
		outputDir:    outputDir,
	}
}

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
		s.log.Error("Failed to retrieve template", err, "template_id", templateID)
		return nil, errors.DatabaseError("failed to retrieve template", err)
	}
	if templateRecord.Layout.FilePath == "" {
		s.log.Error("Layout information missing or invalid for template", nil, "template_id", templateID, "layout_id", templateRecord.LayoutID)
		return nil, errors.InternalServerError("template configuration incomplete: layout file path missing", nil)
	}
	var requiredFields []RequiredFieldConfig
	if err := json.Unmarshal(templateRecord.RequiredFields, &requiredFields); err != nil {
		s.log.Error("Failed to parse required_fields JSON from template", err, "template_id", templateID)
		return nil, errors.InternalServerError("template configuration error: invalid required fields", err)
	}

	validationErrors := make(map[string]string)
	for _, fieldConfig := range requiredFields {
		fieldName := fieldConfig.Name
		userValue, exists := input.Data[fieldName]

		if !exists || userValue == nil || fmt.Sprintf("%v", userValue) == "" {
			// Basic required check (can be enhanced if templates define optional fields)
			validationErrors[fieldName] = fmt.Sprintf("%s is required.", fieldConfig.Label)
			continue // Skip further checks if missing
		}

		// Ensure value is treated as a string for validation checks
		userValueStr := fmt.Sprintf("%v", userValue)

		// MaxLength Check (using rune count for UTF-8 safety)
		if fieldConfig.MaxLength > 0 && utf8.RuneCountInString(userValueStr) > fieldConfig.MaxLength {
			validationErrors[fieldName] = fmt.Sprintf("%s cannot exceed %d characters.", fieldConfig.Label, fieldConfig.MaxLength)
		}

		// Pattern Check
		if fieldConfig.Pattern != "" {
			matched, _ := regexp.MatchString(fieldConfig.Pattern, userValueStr)
			if !matched {
				errorMsg := "Invalid format."
				if fieldConfig.PatternTitle != "" {
					errorMsg = fieldConfig.PatternTitle
				}
				validationErrors[fieldName] = fmt.Sprintf("%s: %s", fieldConfig.Label, errorMsg)
			}
		}
	}

	// If any validation errors occurred, return them immediately
	if len(validationErrors) > 0 {
		s.log.Warn("Backend validation failed for poster input data", validationErrors)
		// Convert map[string]string to map[string]interface{} for ValidationError
		errorDetails := make(map[string]interface{}, len(validationErrors))
		for k, v := range validationErrors {
			errorDetails[k] = v
		}
		return nil, errors.ValidationError("invalid input data provided", nil, errorDetails)
	}
	finalTemplateData := make(map[string]interface{})

	var baseCustomization map[string]interface{}
	if len(templateRecord.DefaultCustomization) > 0 && string(templateRecord.DefaultCustomization) != "null" {
		if err := json.Unmarshal(templateRecord.DefaultCustomization, &baseCustomization); err != nil {
			var strData string
			if err2 := json.Unmarshal(templateRecord.DefaultCustomization, &strData); err2 == nil && strData != "" {
				if err3 := json.Unmarshal([]byte(strData), &baseCustomization); err3 != nil {
					s.log.Error("Failed to unmarshal nested base customization data", err3, "raw_data", strData)
					return nil, errors.InternalServerError("invalid template base customization format", err3)
				}
			} else {
				s.log.Error("Failed to unmarshal base customization data", err, "raw_data", string(templateRecord.DefaultCustomization))
				return nil, errors.InternalServerError("invalid template base customization", err)
			}
		}
	}
	for key, value := range baseCustomization {
		finalTemplateData[key] = value
	}

	if input.CustomizationData != nil {
		for key, value := range input.CustomizationData {
			finalTemplateData[key] = value
		}
	}

	var logoSVG template.HTML = ""
	if logoAssetIDVal, ok := finalTemplateData["header_logo_asset_id"]; ok {
		var logoAssetID uint = 0
		if idFloat, ok := logoAssetIDVal.(float64); ok && idFloat > 0 {
			logoAssetID = uint(idFloat)
		} else if idStr, ok := logoAssetIDVal.(string); ok {
			idUint64, _ := strconv.ParseUint(idStr, 10, 64)
			logoAssetID = uint(idUint64)
		}

		if logoAssetID > 0 {
			asset, err := s.assetRepo.GetAssetByID(ctx, logoAssetID)
			if err == nil && asset != nil && asset.Type == "logo" {
				logoSVG = template.HTML(asset.Data)
				if _, userSetColor := input.CustomizationData["primary_color"]; !userSetColor && asset.DefaultColor != "" {
					finalTemplateData["primary_color"] = asset.DefaultColor
				}
			} else if err != nil && err != gorm.ErrRecordNotFound {
				s.log.Warn("Failed to fetch logo asset", err, "asset_id", logoAssetID)
			} else {
				s.log.Warn("Logo asset not found or not of type 'logo'", nil, "asset_id", logoAssetID)
			}
		}
	}
	finalTemplateData["header_logo_svg"] = logoSVG

	if input.Data != nil {
		for key, value := range input.Data {
			finalTemplateData[key] = value
			if strings.HasSuffix(key, "_number") {
				if strValue, ok := value.(string); ok && strValue != "" {
					finalTemplateData[key+"Split"] = strings.Split(strValue, "")
				} else {
					finalTemplateData[key+"Split"] = []string{}
				}
			}
		}
	}

	finalTemplateData["business_name"] = input.BusinessName

	htmlContent, err := s.renderHTMLTemplate(finalTemplateData, templateRecord.Layout.FilePath)
	if err != nil {
		return nil, errors.InternalServerError("failed to render template", err)
	}

	pdfPath, err := s.renderToPDF(ctx, htmlContent, input.BusinessName, make(map[string]interface{}))
	if err != nil {
		return nil, errors.InternalServerError("failed to generate PDF", err)
	}

	userInputDataJSON, err := json.Marshal(input.Data)
	if err != nil {
		return nil, errors.InternalServerError("failed to marshal user input data", err)
	}
	finalCustomizationJSON, err := json.Marshal(finalTemplateData)
	if err != nil {
		s.log.Error("Failed to marshal final customization data", err, "data", finalTemplateData)
		return nil, errors.InternalServerError("failed to marshal final customization data", err)
	}

	poster := &models.Poster{
		PosterTemplateID:   templateID,
		BusinessName:       input.BusinessName,
		UserInputData:      datatypes.JSON(userInputDataJSON),
		FinalCustomization: datatypes.JSON(finalCustomizationJSON),
		PDFURL:             pdfPath,
		Status:             "completed",
	}

	if err := s.repo.CreatePoster(ctx, poster); err != nil {
		s.log.Error("Failed to save poster to database", err)
		return nil, errors.DatabaseError("failed to save poster", err)
	}

	return &dto.PosterResponse{
		ID:           poster.ID,
		TemplateID:   poster.PosterTemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}, nil
}

func (s *posterSubService) renderHTMLTemplate(data map[string]interface{}, layoutFilePath string) (string, error) {
	templatePath := filepath.Join(s.templatesDir, layoutFilePath)
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		s.log.Error("Failed to read template file", err, "path", templatePath)
		return "", fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}
	tmpl, err := template.New(filepath.Base(layoutFilePath)).Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(string(templateBytes))
	if err != nil {
		s.log.Error("Failed to parse HTML template", err, "path", templatePath)
		return "", fmt.Errorf("failed to parse template %s: %w", layoutFilePath, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		s.log.Error("Failed to execute HTML template", err, "path", templatePath)
		return "", fmt.Errorf("failed to execute template %s: %w", layoutFilePath, err)
	}
	return buf.String(), nil
}

func (s *posterSubService) renderToPDF(ctx context.Context, htmlContent string, businessName string, templateData map[string]interface{}) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var pdfBuffer []byte
	safeBusinessName := strings.ReplaceAll(businessName, " ", "_")
	pdfPath := filepath.Join(s.outputDir, fmt.Sprintf("%s_%d.pdf", safeBusinessName, time.Now().Unix()))
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
			printParams := page.PrintToPDF().
				WithPreferCSSPageSize(true).
				WithPrintBackground(true).
				WithMarginTop(0).WithMarginBottom(0).WithMarginLeft(0).WithMarginRight(0)

			buf, _, err := printParams.Do(ctx)
			if err != nil {
				return err
			}
			pdfBuffer = buf
			return nil
		}),
	)
	if err != nil {
		s.log.Error("Chromedp PDF generation failed", err)
		return "", fmt.Errorf("chromedp failed: %w", err)
	}
	if err := os.WriteFile(pdfPath, pdfBuffer, 0644); err != nil {
		s.log.Error("Failed to write PDF file", err, "path", pdfPath)
		return "", fmt.Errorf("failed to write PDF file: %w", err)
	}
	s.log.Info("PDF generated successfully using CSS @page", "path", pdfPath)
	return pdfPath, nil
}

// GetPosterByID uses the correct PosterRepository interface.
func (s *posterSubService) GetPosterByID(ctx context.Context, id uint) (*dto.PosterResponse, error) {
	s.log.Info("Getting poster by ID", "poster_id", id)
	poster, err := s.repo.GetPosterByID(ctx, id) // Use s.repo (PosterRepository)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Warn("Poster not found", "poster_id", id)
			return nil, errors.NotFoundError("poster not found", err)
		}
		s.log.Error("Failed to get poster by ID", err, "poster_id", id)
		return nil, errors.DatabaseError("failed to retrieve poster", err)
	}
	return &dto.PosterResponse{
		ID:           poster.ID,
		TemplateID:   poster.PosterTemplateID,
		BusinessName: poster.BusinessName,
		PDFURL:       poster.PDFURL,
		Status:       poster.Status,
	}, nil
}
