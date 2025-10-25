package handlers

import (
	"evernos-api2/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// GetAllCategories handles GET /category (PUBLIC)
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch categories",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Categories fetched successfully",
		"data":    categories,
	})
}

// GetCategoryByID handles GET /category/:id (PUBLIC)
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid category ID",
		})
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category fetched successfully",
		"data":    category,
	})
}

// CreateCategory handles POST /category (ADMIN ONLY)
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var request struct {
		NamaCategory string `json:"nama_category"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validasi field yang diperlukan
	if request.NamaCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "nama_category is required",
		})
	}

	category, err := h.categoryService.CreateCategory(request.NamaCategory)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create category",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
		"data":    category,
	})
}

// UpdateCategory handles PUT /category/:id (ADMIN ONLY)
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid category ID",
		})
	}

	var request struct {
		NamaCategory string `json:"nama_category"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validasi field yang diperlukan
	if request.NamaCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "nama_category is required",
		})
	}

	category, err := h.categoryService.UpdateCategory(uint(id), request.NamaCategory)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update category",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category updated successfully",
		"data":    category,
	})
}

// DeleteCategory handles DELETE /category/:id (ADMIN ONLY)
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid category ID",
		})
	}

	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to delete category",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
