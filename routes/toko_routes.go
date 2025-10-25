package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTokoRoutes(app *fiber.App, tokoHandler *handlers.TokoHandler) {
	// Group untuk toko routes
	toko := app.Group("/toko")

	// Public routes (tidak perlu auth)
	toko.Get("/", tokoHandler.GetAllTokos)        // GET /toko (list dengan query params)

	// Protected routes (perlu auth) - harus didefinisikan sebelum route dengan parameter
	toko.Post("/", middleware.AuthMiddleware, tokoHandler.CreateToko)          // POST /toko
	toko.Get("/my", middleware.AuthMiddleware, tokoHandler.GetMyToko)           // GET /toko/my

	// Routes dengan parameter harus didefinisikan terakhir
	toko.Get("/:id_toko", tokoHandler.GetTokoByID) // GET /toko/:id_toko
	toko.Put("/:id_toko", middleware.AuthMiddleware, tokoHandler.UpdateToko)  // PUT /toko/:id_toko
}