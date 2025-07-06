package port

import (
	"context"
	"github.com/envercigal/golang/internal/core/domain"
	"io"
)

type DriverLocationRepository interface {
	Create(ctx context.Context, driverLocations *domain.DriverLocation) (*domain.DriverLocation, error)
	BulkCreate(ctx context.Context, driverLocations []*domain.DriverLocation) error
	FindNearest(ctx context.Context, longitude, latitude float64) (*domain.DriverLocation, error)
}

type DriverLocationService interface {
	Create(ctx context.Context, driverLocations *domain.DriverLocation) (*domain.DriverLocation, error)
	BulkCreate(ctx context.Context, reader io.Reader) error
	FindNearest(ctx context.Context, longitude, latitude float64) (*domain.DriverLocation, error)
}
