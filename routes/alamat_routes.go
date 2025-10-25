package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAlamatRoutes(app *fiber.App, alamatHandler *handlers.AlamatHandler) {
	// Group untuk user alamat dengan authentication middleware
	userGroup := app.Group("/user", middleware.AuthMiddleware)

	// Alamat routes - semua endpoint memerlukan authentication
	alamatGroup := userGroup.Group("/alamat")

	// GET /user/alamat - Ambil semua alamat user
	alamatGroup.Get("/", alamatHandler.GetUserAlamats)

	// GET /user/alamat/:id - Ambil alamat berdasarkan ID
	alamatGroup.Get("/:id", alamatHandler.GetAlamatByID)

	// POST /user/alamat - Buat alamat baru
	alamatGroup.Post("/", alamatHandler.CreateAlamat)

	// PUT /user/alamat/:id - Update alamat
	alamatGroup.Put("/:id", alamatHandler.UpdateAlamat)

	// DELETE /user/alamat/:id - Hapus alamat
	alamatGroup.Delete("/:id", alamatHandler.DeleteAlamat)
}
