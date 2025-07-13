package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"url-shortener/internal/core/domain"
)

type QRHandlerTestSuite struct {
	suite.Suite
	handler       *QRHandler
	mockQRService *MockQRService
}

func TestQRHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(QRHandlerTestSuite))
}

func (suite *QRHandlerTestSuite) SetupTest() {
	suite.mockQRService = &MockQRService{}
	suite.handler = NewQRHandler(suite.mockQRService)
}

func (suite *QRHandlerTestSuite) TestGenerateQRCode_Success() {
	// Setup
	req := domain.QRCodeRequest{
		URL:             "https://example.com",
		Size:            300,
		ErrorCorrection: "M",
	}
	
	qrResponse := &domain.QRCodeResponse{
		Format:   "png",
		Data:     []byte("fake-qr-data"),
		Size:     300,
		URL:      "https://example.com",
		MimeType: "image/png",
	}

	suite.mockQRService.On("GenerateQRCode", mock.Anything, req).Return(qrResponse, nil)

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/qr/generate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCode(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), qrResponse.MimeType, rr.Header().Get("Content-Type"))
	assert.Equal(suite.T(), qrResponse.Format, rr.Header().Get("X-QR-Format"))
	assert.Equal(suite.T(), qrResponse.URL, rr.Header().Get("X-QR-URL"))
	assert.Equal(suite.T(), string(qrResponse.Data), rr.Body.String())
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGenerateQRCode_InvalidRequest() {
	// Create request with invalid JSON
	httpReq := httptest.NewRequest("POST", "/api/qr/generate", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCode(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *QRHandlerTestSuite) TestGenerateQRCode_InvalidOptions() {
	// Setup
	req := domain.QRCodeRequest{
		URL:             "", // Invalid empty URL
		Size:            300,
		ErrorCorrection: "M",
	}

	// Create request
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/qr/generate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCode(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *QRHandlerTestSuite) TestGenerateQRCodeForURL_PNGFormat() {
	// Setup
	shortCode := "test123"
	
	qrResponse := &domain.QRCodeResponse{
		Format:   "png",
		Data:     []byte("fake-png-data"),
		Size:     256,
		URL:      "https://short.ly/" + shortCode,
		MimeType: "image/png",
	}

	suite.mockQRService.On("GenerateQRCodeForURL", mock.Anything, shortCode, mock.AnythingOfType("domain.QRCodeOptions")).Return(qrResponse, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/qr/test123.png?size=256", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortCode", shortCode)
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCodeForURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "image/png", rr.Header().Get("Content-Type"))
	assert.Equal(suite.T(), "public, max-age=3600", rr.Header().Get("Cache-Control"))
	assert.Equal(suite.T(), string(qrResponse.Data), rr.Body.String())
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGenerateQRCodeForURL_SVGFormat() {
	// Setup
	shortCode := "test123"
	
	qrResponse := &domain.QRCodeResponse{
		Format:   "svg",
		Data:     []byte("<svg>fake-svg-data</svg>"),
		Size:     256,
		URL:      "https://short.ly/" + shortCode,
		MimeType: "image/svg+xml",
	}

	suite.mockQRService.On("GenerateQRCodeForURL", mock.Anything, shortCode, mock.AnythingOfType("domain.QRCodeOptions")).Return(qrResponse, nil)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/qr/test123.svg", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortCode", shortCode)
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCodeForURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "image/svg+xml", rr.Header().Get("Content-Type"))
	assert.Equal(suite.T(), string(qrResponse.Data), rr.Body.String())
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGenerateQRCodeForURL_NotFound() {
	// Setup
	shortCode := "invalid"
	
	suite.mockQRService.On("GenerateQRCodeForURL", mock.Anything, shortCode, mock.AnythingOfType("domain.QRCodeOptions")).Return(nil, domain.ErrURLNotFound)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/qr/invalid.png", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortCode", shortCode)
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCodeForURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, rr.Code)
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGenerateQRCodeForURL_InvalidSize() {
	// Setup mock to return error for invalid options
	suite.mockQRService.On("GenerateQRCodeForURL", mock.Anything, "test123", mock.AnythingOfType("domain.QRCodeOptions")).Return(nil, domain.ErrURLNotFound)
	
	// Create request with invalid size
	httpReq := httptest.NewRequest("GET", "/api/qr/test123.png?size=invalid", nil)
	
	// Add chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortCode", "test123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GenerateQRCodeForURL(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, rr.Code)
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGetQRCodeFormats_Success() {
	// Setup
	formats := []string{"png", "svg", "jpeg"}
	
	suite.mockQRService.On("GetQRCodeFormats", mock.Anything).Return(formats)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/qr/formats", nil)
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetQRCodeFormats(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response struct {
		Formats []string `json:"formats"`
	}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), formats, response.Formats)
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGetQRCodeSizes_Success() {
	// Setup
	sizes := []int{128, 256, 512, 1024}
	
	suite.mockQRService.On("GetQRCodeSizes", mock.Anything).Return(sizes)

	// Create request
	httpReq := httptest.NewRequest("GET", "/api/qr/sizes", nil)
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetQRCodeSizes(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response struct {
		Sizes []int `json:"sizes"`
	}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), sizes, response.Sizes)
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestGetQRCodePreview_Success() {
	// Setup
	req := domain.QRCodeRequest{
		URL:             "https://example.com",
		Size:            300,
		Format:          "png",
		ErrorCorrection: "M",
	}
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/qr/preview", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.GetQRCodePreview(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "https://example.com", response["url"])
	assert.Equal(suite.T(), float64(300), response["size"])}

func (suite *QRHandlerTestSuite) TestValidateQRCodeOptions_Success() {
	// Setup
	options := domain.QRCodeOptions{
		Size:             256,
		ErrorCorrection:  "M",
		Border:           4,
		ForegroundColor:  "#000000",
		BackgroundColor:  "#FFFFFF",
	}
	
	suite.mockQRService.On("ValidateQRCodeOptions", mock.Anything, options).Return(nil)

	// Create request
	body, _ := json.Marshal(options)
	httpReq := httptest.NewRequest("POST", "/api/qr/validate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.ValidateQRCodeOptions(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response["valid"])
	
	suite.mockQRService.AssertExpectations(suite.T())
}

func (suite *QRHandlerTestSuite) TestValidateQRCodeOptions_InvalidOptions() {
	// Setup
	options := domain.QRCodeOptions{
		Size:            10, // Too small
		ErrorCorrection: "X", // Invalid level
	}
	
	suite.mockQRService.On("ValidateQRCodeOptions", mock.Anything, options).Return(domain.ErrInvalidRequest)

	// Create request
	body, _ := json.Marshal(options)
	httpReq := httptest.NewRequest("POST", "/api/qr/validate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.ValidateQRCodeOptions(rr, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	
	suite.mockQRService.AssertExpectations(suite.T())
}

// Mock QR service implementation
type MockQRService struct {
	mock.Mock
}

func (m *MockQRService) GenerateQRCode(ctx context.Context, req domain.QRCodeRequest) (*domain.QRCodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.QRCodeResponse), args.Error(1)
}

func (m *MockQRService) GenerateQRCodeForURL(ctx context.Context, shortCode string, options domain.QRCodeOptions) (*domain.QRCodeResponse, error) {
	args := m.Called(ctx, shortCode, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.QRCodeResponse), args.Error(1)
}

func (m *MockQRService) GetQRCodeFormats(ctx context.Context) []string {
	args := m.Called(ctx)
	return args.Get(0).([]string)
}

func (m *MockQRService) GetQRCodeSizes(ctx context.Context) []int {
	args := m.Called(ctx)
	return args.Get(0).([]int)
}

func (m *MockQRService) ValidateQRCodeOptions(ctx context.Context, options domain.QRCodeOptions) error {
	args := m.Called(ctx, options)
	return args.Error(0)
}