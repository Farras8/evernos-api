package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(app *fiber.App, productHandler *handlers.ProductHandler) {
	// Public routes - tidak memerlukan autentikasi
	app.Get("/product", productHandler.GetAllProducts)
	app.Get("/product/:id", productHandler.GetProductByID)

	// Protected routes - memerlukan autentikasi
	app.Post("/product", middleware.AuthMiddleware, productHandler.CreateProduct)
	app.Put("/product/:id", middleware.AuthMiddleware, productHandler.UpdateProduct)
	app.Delete("/product/:id", middleware.AuthMiddleware, productHandler.DeleteProduct)
}
