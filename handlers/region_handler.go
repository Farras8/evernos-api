package handlers

import (
	"evernos-api2/services"

	"github.com/gofiber/fiber/v2"
)

// GetProvinces handles GET /api/provinces
func GetProvinces(c *fiber.Ctx) error {
	provinces, err := services.GetProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch provinces",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Provinces fetched successfully",
		"data":    provinces,
	})
}

// GetCitiesByProvince handles GET /provcity/listcities/:prov_id
func GetCitiesByProvince(c *fiber.Ctx) error {
	provinceID := c.Params("prov_id")
	if provinceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Province ID is required",
		})
	}

	cities, err := services.GetCitiesByProvinceID(provinceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch cities",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cities fetched successfully",
		"data":    cities,
	})
}

// GetAllCities handles GET /api/cities
func GetAllCities(c *fiber.Ctx) error {
	cities, err := services.GetAllCities()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch all cities",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "All cities fetched successfully",
		"data":    cities,
	})
}

// GetProvinceDetail handles GET /provcity/detailprovince/:prov_id
func GetProvinceDetail(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	if provID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Province ID is required",
		})
	}

	province, err := services.GetProvinceByID(provID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Province not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Province detail fetched successfully",
		"data":    province,
	})
}

// GetCityDetail handles GET /provcity/detailcity/:city_id
func GetCityDetail(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	if cityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "City ID is required",
		})
	}

	city, err := services.GetCityByID(cityID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "City not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "City detail fetched successfully",
		"data":    city,
	})
}