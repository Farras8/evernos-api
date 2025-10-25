// file: handlers/auth_handler.go

package handlers

import (
	"evernos-api2/database"
	"evernos-api2/models"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// Validasi field yang wajib diisi
	requiredFields := []string{"password", "email", "nama", "noTelp", "tanggalLahir", "jenisKelamin", "pekerjaan", "idProvinsi", "idKota"}
	for _, field := range requiredFields {
		if data[field] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Field %s is required", field),
			})
		}
	}

	// Parse tanggal lahir
	tanggalLahir, err := time.Parse("2006-01-02", data["tanggalLahir"])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format for tanggalLahir. Use YYYY-MM-DD format",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	// Buat user baru dengan semua field
	user := models.User{
		Nama:         data["nama"],
		Email:        data["email"],
		KataSandi:    string(hashedPassword),
		NoTelp:       data["noTelp"],
		TanggalLahir: tanggalLahir,
		JenisKelamin: data["jenisKelamin"],
		Tentang:      data["tentang"], // Optional field
		Pekerjaan:    data["pekerjaan"],
		IdProvinsi:   data["idProvinsi"],
		IdKota:       data["idKota"],
		IsAdmin:      false, // Default user adalah bukan admin
	}

	// Mulai transaksi database
	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to start transaction"})
	}

	// Buat user
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create user"})
	}

	// Buat toko otomatis setelah user berhasil dibuat
	toko := models.Toko{
		IdUser:   user.ID,
		NamaToko: fmt.Sprintf("Toko %s", user.Nama),
		UrlToko:  fmt.Sprintf("toko-%d", user.ID),
	}

	if err := tx.Create(&toko).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create store"})
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to commit transaction"})
	}

	// Return response tanpa password
	response := fiber.Map{
		"message": "User and store created successfully",
		"user": fiber.Map{
			"id":           user.ID,
			"nama":         user.Nama,
			"email":        user.Email,
			"noTelp":       user.NoTelp,
			"tanggalLahir": user.TanggalLahir.Format("2006-01-02"),
			"jenisKelamin": user.JenisKelamin,
			"tentang":      user.Tentang,
			"pekerjaan":    user.Pekerjaan,
			"idProvinsi":   user.IdProvinsi,
			"idKota":       user.IdKota,
			"isAdmin":      user.IsAdmin,
			"createdAt":    user.CreatedAt,
		},
		"toko": fiber.Map{
			"id":       toko.ID,
			"namaToko": toko.NamaToko,
			"urlToko":  toko.UrlToko,
		},
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Membuat claims untuk JWT
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Membuat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}
