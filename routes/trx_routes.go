package routes

import (
	"evernos-api2/handlers"
	"evernos-api2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTrxRoutes(app *fiber.App, trxHandler *handlers.TrxHandler) {
	// Semua endpoint transaksi memerlukan autentikasi
	trx := app.Group("/trx", middleware.AuthMiddleware)

	// GET /trx - Mengambil semua transaksi user dengan pagination
	trx.Get("/", trxHandler.GetAllTrx)

	// GET /trx/:id - Mengambil transaksi berdasarkan ID
	trx.Get("/:id", trxHandler.GetTrxByID)

	// POST /trx - Membuat transaksi baru
	trx.Post("/", trxHandler.CreateTrx)
}