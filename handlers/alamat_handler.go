package handlers

import (
	"evernos-api2/models"
	"evernos-api2/services"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlamatHandler struct {
	alamatService *services.AlamatService
}

func NewAlamatHandler(alamatService *services.AlamatService) *AlamatHandler {
	return &AlamatHandler{alamatService: alamatService}
}

// GetUserAlamats mengambil semua alamat user yang sedang login
func (h *AlamatHandler) GetUserAlamats(c *fiber.Ctx) error {
	// Ambil user ID dari context (dari middleware auth)
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	userID := uint(userIDFloat)

	alamats, err := h.alamatService.GetUserAlamats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data alamat",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data alamat berhasil diambil",
		"data":    alamats,
	})
}

// GetAlamatByID mengambil alamat berdasarkan ID
func (h *AlamatHandler) GetAlamatByID(c *fiber.Ctx) error {
	// Ambil user ID dari context
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	userID := uint(userIDFloat)

	// Ambil ID dari parameter
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID alamat tidak valid",
		})
	}

	alamat, err := h.alamatService.GetAlamatByID(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data alamat berhasil diambil",
		"data":    alamat,
	})
}

// CreateAlamat membuat alamat baru
func (h *AlamatHandler) CreateAlamat(c *fiber.Ctx) error {
	// Ambil user ID dari context
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	userID := uint(userIDFloat)

	var alamat models.Alamat
	if err := c.BodyParser(&alamat); err != nil {
		// Log error untuk debugging
		fmt.Printf("BodyParser error: %v\n", err)
		fmt.Printf("Request body: %s\n", string(c.Body()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
			"detail": err.Error(),
		})
	}

	// Set user ID
	alamat.IdUser = userID

	if err := h.alamatService.CreateAlamat(&alamat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alamat berhasil dibuat",
		"data":    alamat,
	})
}

// UpdateAlamat memperbarui alamat
func (h *AlamatHandler) UpdateAlamat(c *fiber.Ctx) error {
	// Ambil user ID dari context
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	userID := uint(userIDFloat)

	// Ambil ID dari parameter
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID alamat tidak valid",
		})
	}

	var alamat models.Alamat
	if err := c.BodyParser(&alamat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	if err := h.alamatService.UpdateAlamat(uint(id), userID, &alamat); err != nil {
		if err.Error() == "alamat tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Alamat berhasil diperbarui",
		"data":    alamat,
	})
}

// DeleteAlamat menghapus alamat
func (h *AlamatHandler) DeleteAlamat(c *fiber.Ctx) error {
	// Ambil user ID dari context
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	userID := uint(userIDFloat)

	// Ambil ID dari parameter
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID alamat tidak valid",
		})
	}

	if err := h.alamatService.DeleteAlamat(uint(id), userID); err != nil {
		if err.Error() == "alamat tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus alamat",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Alamat berhasil dihapus",
	})
}