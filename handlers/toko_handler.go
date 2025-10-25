package handlers

import (
	"evernos-api2/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TokoHandler struct {
	tokoService *services.TokoService
}

func NewTokoHandler(tokoService *services.TokoService) *TokoHandler {
	return &TokoHandler{tokoService: tokoService}
}

// CreateToko membuat toko baru untuk user
func (h *TokoHandler) CreateToko(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Parse request body
	var tokoData map[string]interface{}
	if err := c.BodyParser(&tokoData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	toko, err := h.tokoService.CreateToko(uint(userID), tokoData)
	if err != nil {
		if err.Error() == "user sudah memiliki toko" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil membuat toko",
		"data":    toko,
	})
}

// GetMyToko mengambil toko milik user yang sedang login
func (h *TokoHandler) GetMyToko(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	toko, err := h.tokoService.GetMyToko(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data toko",
		"data":    toko,
	})
}

// GetTokoByID mengambil toko berdasarkan ID
func (h *TokoHandler) GetTokoByID(c *fiber.Ctx) error {
	idParam := c.Params("id_toko")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID toko tidak valid",
		})
	}

	toko, err := h.tokoService.GetTokoByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data toko",
		"data":    toko,
	})
}

// UpdateToko memperbarui toko
func (h *TokoHandler) UpdateToko(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	idParam := c.Params("id_toko")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID toko tidak valid",
		})
	}

	// Parse request body
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	toko, err := h.tokoService.UpdateToko(uint(id), uint(userID), updateData)
	if err != nil {
		if err.Error() == "anda tidak memiliki akses untuk mengupdate toko ini" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengupdate toko",
		"data":    toko,
	})
}

// GetAllTokos mengambil semua toko dengan pagination dan filter
func (h *TokoHandler) GetAllTokos(c *fiber.Ctx) error {
	// Ambil query parameters
	limit := c.Query("limit")
	page := c.Query("page")
	namaToko := c.Query("nama_toko")

	tokos, pagination, err := h.tokoService.GetAllTokos(limit, page, namaToko)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Berhasil mengambil data toko",
		"data":       tokos,
		"pagination": pagination,
	})
}