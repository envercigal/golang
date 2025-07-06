package http

import (
	"context"
	"errors"
	"github.com/envercigal/golang/internal/core/domain"
	"github.com/envercigal/golang/internal/core/port"
	circuitbreaker "github.com/envercigal/golang/pkg"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

func CreateDriverHandler(svc port.DriverLocationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dl domain.DriverLocation
		if err := c.BodyParser(&dl); err != nil {
			return fiber.ErrBadRequest
		}
		created, err := svc.Create(c.Context(), &dl)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}
}

func ImportDriversHandler(svc port.DriverLocationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return fiber.ErrBadRequest
		}
		f, err := file.Open()

		if err != nil {
			return fiber.ErrInternalServerError
		}

		bulkErr := svc.BulkCreate(context.Background(), f)

		if bulkErr != nil {
			return fiber.ErrInternalServerError
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "import success"})
	}
}

func FindNearestHandler(svc port.DriverLocationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lon, err1 := strconv.ParseFloat(c.Query("lon"), 64)
		lat, err2 := strconv.ParseFloat(c.Query("lat"), 64)
		if err1 != nil || err2 != nil {
			return fiber.ErrBadRequest
		}

		driver, err := svc.FindNearest(c.Context(), lon, lat)

		// Circuit breaker logic can be added here if needed.
		// For testing purposes, only the HTTP status codes were changed.
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return fiber.ErrNotFound
		case errors.Is(err, circuitbreaker.ErrOpen):
			return fiber.ErrServiceUnavailable
		case errors.Is(err, circuitbreaker.ErrHalfOpen):
			return fiber.ErrTooManyRequests
		case err != nil:
			return fiber.ErrInternalServerError
		default:
			return c.Status(http.StatusOK).JSON(driver)
		}
	}
}
