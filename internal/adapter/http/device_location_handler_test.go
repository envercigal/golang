package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"golang-case/internal/core/domain"
	"golang-case/internal/core/port"
)

// fake svc
type fakeService struct {
	createFn      func(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error)
	bulkCreateFn  func(ctx context.Context, r io.Reader) error
	findNearestFn func(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error)
}

func (f *fakeService) Create(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error) {
	return f.createFn(ctx, dl)
}

func (f *fakeService) BulkCreate(ctx context.Context, r io.Reader) error {
	return f.bulkCreateFn(ctx, r)
}

func (f *fakeService) FindNearest(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
	return f.findNearestFn(ctx, lon, lat)
}

func setupApp(svc port.DriverLocationService) *fiber.App {
	app := fiber.New()
	RegisterDriverRoutes(app, svc)
	return app
}

func TestCreateHandler(t *testing.T) {
	called := false
	svc := &fakeService{
		createFn: func(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error) {
			called = true
			dl.DriverID = 7
			return dl, nil
		},
	}
	app := setupApp(svc)

	body := `{"driver_id":5,"location":{"type":"Point","coordinates":[29,41]}}`
	req := httptest.NewRequest("POST", "/drivers/", strings.NewReader(body))
	token := makeTestToken()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.True(t, called)
}

func TestImportHandler(t *testing.T) {
	done := false
	svc := &fakeService{
		bulkCreateFn: func(ctx context.Context, r io.Reader) error {
			done = true
			return nil
		},
	}
	app := setupApp(svc)

	// prepare multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	_, writeErr := part.Write([]byte("lat,lon\n41,29\n"))
	if writeErr != nil {
		return
	}
	err := writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest("POST", "/drivers/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	token := makeTestToken()
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.True(t, done)
}

func makeTestToken() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"authenticated":true}`))
	return header + "." + payload + "."
}

func TestFindNearestHandler(t *testing.T) {
	svc := &fakeService{
		findNearestFn: func(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
			return &domain.DriverLocation{DriverID: 9}, nil
		},
	}
	app := setupApp(svc)

	req := httptest.NewRequest("GET", "/drivers/nearest?lon=29&lat=41", nil)
	token := makeTestToken()
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
