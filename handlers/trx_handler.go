package handlers

import (
	"evernos-api2/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TrxHandler struct {
	trxService *services.TrxService
}

func NewTrxHandler(trxService *services.TrxService) *TrxHandler {
	return &TrxHandler{trxService: trxService}
}

// GetAllTrx mengambil semua transaksi user dengan pagination
func (h *TrxHandler) GetAllTrx(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Ambil query parameters
	limit := c.Query("limit")
	page := c.Query("page")

	trxs, pagination, err := h.trxService.GetAllTrx(uint(userID), limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Berhasil mengambil data transaksi",
		"data":       trxs,
		"pagination": pagination,
	})
}

// GetTrxByID mengambil transaksi berdasarkan ID
func (h *TrxHandler) GetTrxByID(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID transaksi tidak valid",
		})
	}

	trx, err := h.trxService.GetTrxByID(uint(id), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data transaksi",
		"data":    trx,
	})
}

// CreateTrx membuat transaksi baru
func (h *TrxHandler) CreateTrx(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Parse request body
	var trxData map[string]interface{}
	if err := c.BodyParser(&trxData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	trx, err := h.trxService.CreateTrxWithResponse(uint(userID), trxData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil membuat transaksi",
		"data":    trx,
	})
}