package service

import (
	"context"
	"errors"
	"github.com/envercigal/golang/internal/core/domain"
	circuitbreaker "github.com/envercigal/golang/pkg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockRepo implements port.DriverLocationRepository
type mockRepo struct {
	createFn      func(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error)
	bulkCreateFn  func(ctx context.Context, dls []*domain.DriverLocation) error
	findNearestFn func(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error)
}

func (m *mockRepo) Create(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error) {
	return m.createFn(ctx, dl)
}

func (m *mockRepo) BulkCreate(ctx context.Context, dls []*domain.DriverLocation) error {
	return m.bulkCreateFn(ctx, dls)
}

func (m *mockRepo) FindNearest(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
	return m.findNearestFn(ctx, lon, lat)
}

func TestCreate_ValidCoordinates(t *testing.T) {
	now := time.Now().UTC()
	input := &domain.DriverLocation{
		DriverID: 123,
		Location: domain.GeoJSONPoint{Type: "Point", Coordinates: []float64{29.0, 41.0}},
	}

	called := false

	repo := &mockRepo{
		createFn: func(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error) {
			called = true
			// echo back with Updated set
			dl.Updated = now
			return dl, nil
		},
	}
	svc := NewDriverLocationService(repo, circuitbreaker.New(5, 1))
	created, err := svc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, now, created.Updated)
}

func TestCreate_InvalidLat(t *testing.T) {
	input := &domain.DriverLocation{
		DriverID: 1,
		Location: domain.GeoJSONPoint{Type: "Point", Coordinates: []float64{29.0, 200.0}},
	}
	svc := NewDriverLocationService(&mockRepo{}, circuitbreaker.New(5, 1))
	_, err := svc.Create(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "latitude out of range")
}

func TestFindNearest_Success(t *testing.T) {
	expected := &domain.DriverLocation{DriverID: 42}
	repo := &mockRepo{
		findNearestFn: func(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
			return expected, nil
		},
	}
	svc := NewDriverLocationService(repo, circuitbreaker.New(5, 1))
	got, err := svc.FindNearest(context.Background(), 29, 41)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestFindNearest_NotFound(t *testing.T) {
	repo := &mockRepo{
		findNearestFn: func(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
			return nil, errors.New("no documents")
		},
	}
	svc := NewDriverLocationService(repo, circuitbreaker.New(5, 1))
	got, err := svc.FindNearest(context.Background(), 29, 41)
	assert.Error(t, err)
	assert.Nil(t, got)
}
