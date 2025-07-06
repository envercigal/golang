package http

import (
	"github.com/envercigal/golang/internal/adapter/middleware"
	"github.com/envercigal/golang/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

func RegisterDriverRoutes(app *fiber.App, svc port.DriverLocationService) {
	grp := app.Group("/drivers", middleware.RequireAuthenticated())

	grp.Post("/", CreateDriverHandler(svc))
	grp.Post("/import", ImportDriversHandler(svc))
	grp.Get("/nearest", FindNearestHandler(svc))
}
