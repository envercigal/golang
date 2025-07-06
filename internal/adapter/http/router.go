package http

import (
	"github.com/gofiber/fiber/v2"
	"golang-case/internal/adapter/middleware"
	"golang-case/internal/core/port"
)

func RegisterDriverRoutes(app *fiber.App, svc port.DriverLocationService) {
	grp := app.Group("/drivers", middleware.RequireAuthenticated())

	grp.Post("/", CreateDriverHandler(svc))
	grp.Post("/import", ImportDriversHandler(svc))
	grp.Get("/nearest", FindNearestHandler(svc))
}
