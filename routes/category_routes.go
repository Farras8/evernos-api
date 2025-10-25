package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoutes(app *fiber.App, categoryHandler *handlers.CategoryHandler) {
	// Public category routes
	category := app.Group("/category")
	category.Get("/", categoryHandler.GetAllCategories)
	category.Get("/:id", categoryHandler.GetCategoryByID)

	// Admin-only category routes
	adminCategory := app.Group("/category", middleware.AuthMiddleware, middleware.AdminMiddleware)
	adminCategory.Post("/", categoryHandler.CreateCategory)
	adminCategory.Put("/:id", categoryHandler.UpdateCategory)
	adminCategory.Delete("/:id", categoryHandler.DeleteCategory)
}
