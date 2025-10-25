package handlers

import (
	"evernos-api2/services"

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	profileService services.ProfileService
}

func NewProfileHandler(profileService services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile mengambil data profil user yang sedang login
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	// Ambil user_id dari middleware auth
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Convert userID to uint
	// userID dari JWT claims adalah float64, bukan string
	userIDFloat, ok := userID.(float64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID format",
		})
	}
	userIDUint := uint(userIDFloat)

	user, err := h.profileService.GetProfile(uint(userIDUint))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Return user profile tanpa password
	return c.JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"data": fiber.Map{
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
			"updatedAt":    user.UpdatedAt,
		},
	})
}

// UpdateProfile mengupdate data profil user yang sedang login
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	// Ambil user_id dari middleware auth
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Convert userID to uint
	// userID dari JWT claims adalah float64, bukan string
	userIDFloat, ok := userID.(float64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID format",
		})
	}
	userIDUint := uint(userIDFloat)

	user, err := h.profileService.UpdateProfile(uint(userIDUint), data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Return updated profile tanpa password
	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"data": fiber.Map{
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
			"updatedAt":    user.UpdatedAt,
		},
	})
}