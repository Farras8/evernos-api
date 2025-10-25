package database

import (
	"evernos-api2/models"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SeedFirstAdmin creates the first admin user if no admin exists
func SeedFirstAdmin() {
	// Check if admin user with specific email already exists
	var existingAdmin models.User
	err := DB.Where("email = ? OR is_admin = ?", "admin@evernos.com", true).First(&existingAdmin).Error

	if err == nil {
		fmt.Println("👤 Admin user already exists, skipping seeding")
		fmt.Printf("📧 Existing admin email: %s\n", existingAdmin.Email)
		return
	}

	// Hash password for the first admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("❌ Failed to hash admin password: %v\n", err)
		return
	}

	// Create first admin user with unique phone number
	admin := models.User{
		Nama:         "Super Admin",
		Email:        "admin@evernos.com",
		KataSandi:    string(hashedPassword),
		NoTelp:       "081999888777", // Different phone number
		TanggalLahir: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		JenisKelamin: "L",
		Tentang:      "System Administrator",
		Pekerjaan:    "Administrator",
		IdProvinsi:   "11",
		IdKota:       "1101",
		IsAdmin:      true,
	}

	// Start transaction
	tx := DB.Begin()
	if tx.Error != nil {
		fmt.Printf("❌ Failed to start transaction: %v\n", tx.Error)
		return
	}

	// Create admin user
	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		fmt.Printf("❌ Failed to create admin user: %v\n", err)
		return
	}

	// Create admin store
	adminStore := models.Toko{
		IdUser:   admin.ID,
		NamaToko: "Admin Store",
		UrlToko:  "admin-store",
	}

	if err := tx.Create(&adminStore).Error; err != nil {
		tx.Rollback()
		fmt.Printf("❌ Failed to create admin store: %v\n", err)
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("❌ Failed to commit transaction: %v\n", err)
		return
	}

	fmt.Println("✅ First admin user created successfully!")
	fmt.Println("📧 Email: admin@evernos.com")
	fmt.Println("🔑 Password: admin123")
	fmt.Println("👤 Use these credentials to login and get admin token")
}
