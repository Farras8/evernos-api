package handlers

import (
	"evernos-api2/services"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService    *services.ProductService
	fotoProdukService *services.FotoProdukService
	tokoService       *services.TokoService
}

func NewProductHandler(productService *services.ProductService, fotoProdukService *services.FotoProdukService, tokoService *services.TokoService) *ProductHandler {
	return &ProductHandler{
		productService:    productService,
		fotoProdukService: fotoProdukService,
		tokoService:       tokoService,
	}
}

// GetAllProducts mengambil semua produk dengan filtering dan pagination
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	// Ambil query parameters
	filters := map[string]string{
		"nama_produk": c.Query("nama_produk"),
		"limit":       c.Query("limit"),
		"page":        c.Query("page"),
		"category_id": c.Query("category_id"),
		"toko_id":     c.Query("toko_id"),
		"max_harga":   c.Query("max_harga"),
		"min_harga":   c.Query("min_harga"),
	}

	products, pagination, err := h.productService.GetAllProducts(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Berhasil mengambil data produk",
		"data":       products,
		"pagination": pagination,
	})
}

// GetProductByID mengambil produk berdasarkan ID
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID produk tidak valid",
		})
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data produk",
		"data":    product,
	})
}

// CreateProduct membuat produk baru dengan foto
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Parse form data
	productData := make(map[string]interface{})
	
	// Ambil data dari form fields
	if namaProduk := c.FormValue("nama_produk"); namaProduk != "" {
		productData["nama_produk"] = namaProduk
	}
	if hargaReseller := c.FormValue("harga_reseller"); hargaReseller != "" {
		productData["harga_reseller"] = hargaReseller
	}
	if hargaKonsumen := c.FormValue("harga_konsumen"); hargaKonsumen != "" {
		productData["harga_konsumen"] = hargaKonsumen
	}
	if stokStr := c.FormValue("stok"); stokStr != "" {
		if stok, err := strconv.ParseFloat(stokStr, 64); err == nil {
			productData["stok"] = stok
		}
	}
	if deskripsi := c.FormValue("deskripsi"); deskripsi != "" {
		productData["deskripsi"] = deskripsi
	}
	if idCategoryStr := c.FormValue("id_category"); idCategoryStr != "" {
		if idCategory, err := strconv.ParseFloat(idCategoryStr, 64); err == nil {
			productData["id_category"] = idCategory
		}
	}
	if idTokoStr := c.FormValue("id_toko"); idTokoStr != "" {
		if idToko, err := strconv.ParseFloat(idTokoStr, 64); err == nil {
			productData["id_toko"] = idToko
		}
	}

	// Get user's toko ID from database
	tokoID, err := h.tokoService.GetTokoIDByUserID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User belum memiliki toko. Silakan buat toko terlebih dahulu.",
		})
	}
	productData["id_toko"] = float64(tokoID)

	// Handle foto upload (optional)
	var photoFilename string
	file, err := c.FormFile("photo")
	if err == nil && file != nil {
		// Validasi ukuran file (maksimal 5MB)
		if file.Size > 5*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ukuran file terlalu besar. Maksimal 5MB",
			})
		}

		// Validasi tipe file
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".webp"}
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		isValidType := false
		for _, allowedType := range allowedTypes {
			if fileExt == allowedType {
				isValidType = true
				break
			}
		}

		if !isValidType {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tipe file tidak didukung. Gunakan JPG, JPEG, PNG, atau WEBP",
			})
		}

		// Generate nama file unik
		timestamp := time.Now().Format("20060102_150405")
		uniqueID := uuid.New().String()[:8]
		photoFilename = fmt.Sprintf("product_%d_%s_%s%s", uint(userID), timestamp, uniqueID, fileExt)

		// Path untuk menyimpan file
		uploadPath := filepath.Join("uploads", "products", photoFilename)

		// Simpan file
		if err := c.SaveFile(file, uploadPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan file foto",
			})
		}
	}

	// Buat produk
	product, err := h.productService.CreateProduct(uint(userID), productData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Jika ada foto, simpan ke database
	if photoFilename != "" {
		photoURL := fmt.Sprintf("/uploads/products/%s", photoFilename)
		_, err = h.fotoProdukService.AddPhotoToProduct(product.ID, uint(userID), photoURL)
		if err != nil {
			// Log error tapi jangan gagalkan pembuatan produk
			// Product sudah berhasil dibuat, foto gagal disimpan
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil membuat produk",
		"data":    product,
	})
}

// UpdateProduct memperbarui produk dengan foto
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
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
			"error": "ID produk tidak valid",
		})
	}

	// Parse form data
	updateData := make(map[string]interface{})
	
	// Ambil data dari form fields (hanya yang ada)
	if namaProduk := c.FormValue("nama_produk"); namaProduk != "" {
		updateData["nama_produk"] = namaProduk
	}
	if hargaReseller := c.FormValue("harga_reseller"); hargaReseller != "" {
		updateData["harga_reseller"] = hargaReseller
	}
	if hargaKonsumen := c.FormValue("harga_konsumen"); hargaKonsumen != "" {
		updateData["harga_konsumen"] = hargaKonsumen
	}
	if stokStr := c.FormValue("stok"); stokStr != "" {
		if stok, err := strconv.ParseFloat(stokStr, 64); err == nil {
			updateData["stok"] = stok
		}
	}
	if deskripsi := c.FormValue("deskripsi"); deskripsi != "" {
		updateData["deskripsi"] = deskripsi
	}
	if idCategoryStr := c.FormValue("id_category"); idCategoryStr != "" {
		if idCategory, err := strconv.ParseFloat(idCategoryStr, 64); err == nil {
			updateData["id_category"] = idCategory
		}
	}

	// Handle foto upload (optional)
	var photoFilename string
	file, err := c.FormFile("photo")
	if err == nil && file != nil {
		// Validasi ukuran file (maksimal 5MB)
		if file.Size > 5*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ukuran file terlalu besar. Maksimal 5MB",
			})
		}

		// Validasi tipe file
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".webp"}
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		isValidType := false
		for _, allowedType := range allowedTypes {
			if fileExt == allowedType {
				isValidType = true
				break
			}
		}

		if !isValidType {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tipe file tidak didukung. Gunakan JPG, JPEG, PNG, atau WEBP",
			})
		}

		// Generate nama file unik
		timestamp := time.Now().Format("20060102_150405")
		uniqueID := uuid.New().String()[:8]
		photoFilename = fmt.Sprintf("product_%d_%s_%s%s", uint(userID), timestamp, uniqueID, fileExt)

		// Path untuk menyimpan file
		uploadPath := filepath.Join("uploads", "products", photoFilename)

		// Simpan file
		if err := c.SaveFile(file, uploadPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan file foto",
			})
		}
	}

	// Update produk
	product, err := h.productService.UpdateProduct(uint(id), uint(userID), updateData)
	if err != nil {
		if err.Error() == "anda tidak memiliki akses untuk mengupdate produk ini" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Jika ada foto baru, simpan ke database
	if photoFilename != "" {
		photoURL := fmt.Sprintf("/uploads/products/%s", photoFilename)
		_, err = h.fotoProdukService.AddPhotoToProduct(uint(id), uint(userID), photoURL)
		if err != nil {
			// Log error tapi jangan gagalkan update produk
			// Product sudah berhasil diupdate, foto gagal disimpan
		}
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengupdate produk",
		"data":    product,
	})
}

// DeleteProduct menghapus produk
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
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
			"error": "ID produk tidak valid",
		})
	}

	err = h.productService.DeleteProduct(uint(id), uint(userID))
	if err != nil {
		if err.Error() == "anda tidak memiliki akses untuk menghapus produk ini" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "produk tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil menghapus produk",
	})
}