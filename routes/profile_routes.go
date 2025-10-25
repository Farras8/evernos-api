package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupProfileRoutes(app *fiber.App, profileHandler *handlers.ProfileHandler) {
	// Grup untuk rute yang butuh autentikasi (wajib pakai token)
	api := app.Group("/api", middleware.AuthMiddleware)

	// Profile routes
	api.Get("/profile", profileHandler.GetProfile)
	api.Put("/profile", profileHandler.UpdateProfile)
}