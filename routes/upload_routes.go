package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUploadRoutes(app *fiber.App, uploadHandler *handlers.UploadHandler) {
	// Group untuk upload dengan authentication middleware
	upload := app.Group("/upload", middleware.AuthMiddleware)

	// POST /upload/product/assign - Upload dan assign single foto ke produk
	upload.Post("/product/assign", uploadHandler.UploadAndAssignToProduct)

	// POST /upload/product/assign-multiple - Upload dan assign multiple foto ke produk
	upload.Post("/product/assign-multiple", uploadHandler.UploadMultipleAndAssignToProduct)

	// DELETE /product/photo/:foto_id - Hapus foto produk berdasarkan ID foto
	app.Delete("/product/photo/:foto_id", middleware.AuthMiddleware, uploadHandler.DeleteProductPhoto)

	// GET /product/photos/:product_id - Ambil semua foto dari produk tertentu (public)
	app.Get("/product/photos/:product_id", uploadHandler.GetProductPhotos)
}