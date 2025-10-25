// file: main.go

package main

import (
	"log"
	"evernos-api2/database"
	"evernos-api2/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Muat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inisialisasi aplikasi Fiber
	app := fiber.New()

	// Hubungkan & migrasi database
	database.ConnectDB()
	database.MigrateDB()

	// Seed first admin user if none exists
	database.SeedFirstAdmin()

	// Setup semua routes
	routes.SetupRoutes(app)

	// Jalankan server di port 3001
	log.Fatal(app.Listen(":3001"))
}