package routes

import (
	"evernos-api2/database"
	"evernos-api2/handlers"
	"evernos-api2/middleware"
	"evernos-api2/repositories"
	"evernos-api2/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Setup dependencies
	profileRepo := repositories.NewProfileRepository(database.DB)
	profileService := services.NewProfileService(profileRepo)
	profileHandler := handlers.NewProfileHandler(profileService)

	// Category dependencies
	categoryRepo := repositories.NewCategoryRepository(database.DB)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Alamat dependencies
	alamatRepo := repositories.NewAlamatRepository(database.DB)
	alamatService := services.NewAlamatService(alamatRepo)
	alamatHandler := handlers.NewAlamatHandler(alamatService)

	// Toko dependencies
	tokoRepo := repositories.NewTokoRepository(database.DB)
	tokoService := services.NewTokoService(tokoRepo)
	tokoHandler := handlers.NewTokoHandler(tokoService)

	// Product dependencies
	productRepo := repositories.NewProductRepository(database.DB)
	productService := services.NewProductService(productRepo)

	// FotoProduk dependencies
	fotoProdukRepo := repositories.NewFotoProdukRepository(database.DB)
	fotoProdukService := services.NewFotoProdukService(fotoProdukRepo, productRepo)

	// Product handler (needs productService, fotoProdukService, and tokoService)
	productHandler := handlers.NewProductHandler(productService, fotoProdukService, tokoService)

	// LogProduk dependencies
	logProdukRepo := repositories.NewLogProdukRepository(database.DB)
	logProdukService := services.NewLogProdukService(logProdukRepo, productRepo)

	// Trx dependencies
	trxRepo := repositories.NewTrxRepository(database.DB)
	trxService := services.NewTrxService(trxRepo, logProdukService)
	trxHandler := handlers.NewTrxHandler(trxService)

	// Upload dependencies
	uploadHandler := handlers.NewUploadHandler(fotoProdukService)

	// Static file serving untuk uploads
	app.Static("/uploads", "./uploads")

	// Auth routes (public)
	auth := app.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	// Province and City routes (public)
	provcity := app.Group("/provcity")
	provcity.Get("/listprovincies", handlers.GetProvinces)
	provcity.Get("/listcities/:prov_id", handlers.GetCitiesByProvince)
	provcity.Get("/detailprovince/:prov_id", handlers.GetProvinceDetail)
	provcity.Get("/detailcity/:city_id", handlers.GetCityDetail)

	// Category routes (public GET, admin POST/PUT/DELETE)
	SetupCategoryRoutes(app, categoryHandler)

	// Alamat routes (user authentication required)
	SetupAlamatRoutes(app, alamatHandler)

	// Toko routes (mixed public and protected)
	SetupTokoRoutes(app, tokoHandler)

	// Product routes (mixed public and protected)
	SetupProductRoutes(app, productHandler)

	// Trx routes (authentication required)
	SetupTrxRoutes(app, trxHandler)

	// Upload routes (authentication required)
	SetupUploadRoutes(app, uploadHandler)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware)

	// Profile routes
	api.Get("/profile", profileHandler.GetProfile)
	api.Put("/profile", profileHandler.UpdateProfile)

	// Admin routes
	api.Get("/admin/data", middleware.AdminMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "This is a secret admin data.",
		})
	})
	

}