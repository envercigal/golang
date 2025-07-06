package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"golang-case/internal/core/domain"
	"golang-case/internal/core/port"
)

type driverLocationService struct {
	repo       port.DriverLocationRepository
	batchSize  int
	maxWorkers int
}

func NewDriverLocationService(r port.DriverLocationRepository) port.DriverLocationService {
	return &driverLocationService{
		repo:       r,
		batchSize:  10000,
		maxWorkers: 100,
	}
}

func (s *driverLocationService) Create(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error) {
	if len(dl.Location.Coordinates) != 2 {
		return nil, fmt.Errorf("invalid coordinates length: %+v", dl.Location.Coordinates)
	}

	if len(dl.Location.Coordinates) != 2 {
		return nil, fmt.Errorf("invalid coordinates length: %+v", dl.Location.Coordinates)
	}

	lon, lat := dl.Location.Coordinates[0], dl.Location.Coordinates[1]
	if err := validateCoords(lat, lon); err != nil {
		return nil, err
	}

	dl.Updated = time.Now().UTC()
	return s.repo.Create(ctx, dl)
}

func (s *driverLocationService) BulkCreate(ctx context.Context, reader io.Reader) error {
	csvFile := csv.NewReader(reader)
	if _, err := csvFile.Read(); err != nil && err != io.EOF {
		return err
	}

	jobs := make(chan []*domain.DriverLocation, s.maxWorkers)
	var wg sync.WaitGroup

	for i := 0; i < s.maxWorkers; i++ {
		wg.Add(1)
		go s.startWorker(ctx, &wg, jobs)
	}

	if err := s.produceBatches(csvFile, jobs); err != nil {
		close(jobs)
		wg.Wait()
		return err
	}

	close(jobs)
	wg.Wait()
	return nil
}

func (s *driverLocationService) FindNearest(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
	return s.repo.FindNearest(ctx, lon, lat)
}

func (s *driverLocationService) startWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan []*domain.DriverLocation) {
	defer wg.Done()
	for batch := range jobs {
		log.Printf("bulk started batch %+v", batch)
		if err := s.repo.BulkCreate(ctx, batch); err != nil {
			log.Printf("batch import error: %v", err)
		}
	}
}

func (s *driverLocationService) produceBatches(r *csv.Reader, jobs chan<- []*domain.DriverLocation) error {
	now := time.Now().UTC()
	var batch []*domain.DriverLocation
	row := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("csv read error: %v", err)
			continue
		}

		row++
		dl, err := toDriverLocation(record, now, row)
		if err != nil {
			log.Printf("parse record error: %v", err)
			continue
		}

		batch = append(batch, dl)
		if len(batch) >= s.batchSize {
			jobs <- batch
			batch = nil
		}
	}
	if len(batch) > 0 {
		jobs <- batch
	}
	return nil
}

func toDriverLocation(record []string, updated time.Time, row int) (*domain.DriverLocation, error) {
	if len(record) < 2 {
		return nil, errors.New("invalid record length")
	}
	lat, err := strconv.ParseFloat(record[0], 64)

	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}

	if err := validateCoords(lat, lon); err != nil {
		return nil, fmt.Errorf("validation error on row %d: %w", row, err)
	}

	return &domain.DriverLocation{
		DriverID: row,
		Location: domain.GeoJSONPoint{
			Type:        "Point",
			Coordinates: []float64{lon, lat},
		},
		Updated: updated,
	}, nil
}

func validateCoords(lat, lon float64) error {
	if lat < -180 || lat > 180 {
		return fmt.Errorf("latitude out of range: %v", lat)
	}
	if lon < -90 || lon > 90 {
		return fmt.Errorf("longitude out of range: %v", lon)
	}
	return nil
}
