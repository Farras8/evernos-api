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

type UploadHandler struct {
	fotoProdukService *services.FotoProdukService
}

func NewUploadHandler(fotoProdukService *services.FotoProdukService) *UploadHandler {
	return &UploadHandler{fotoProdukService: fotoProdukService}
}



// UploadAndAssignToProduct mengupload foto dan langsung assign ke produk
func (h *UploadHandler) UploadAndAssignToProduct(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Ambil product_id dari form data
	productIDStr := c.FormValue("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id harus disertakan",
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id harus berupa angka",
		})
	}

	// Parse multipart form
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File foto tidak ditemukan. Gunakan field 'photo' untuk upload",
		})
	}

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
	fileName := fmt.Sprintf("product_%d_%s_%s%s", uint(userID), timestamp, uniqueID, fileExt)

	// Path untuk menyimpan file
	uploadPath := filepath.Join("uploads", "products", fileName)

	// Simpan file
	if err := c.SaveFile(file, uploadPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan file",
		})
	}

	// URL file yang bisa diakses
	fileURL := fmt.Sprintf("/uploads/products/%s", fileName)

	// Assign foto ke produk melalui service
	fotoProduk, err := h.fotoProdukService.AddPhotoToProduct(uint(productID), uint(userID), fileURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Foto berhasil diupload dan ditambahkan ke produk",
		"filename":    fileName,
		"url":         fileURL,
		"size":        file.Size,
		"foto_produk": fotoProduk,
	})
}

// UploadMultipleAndAssignToProduct mengupload multiple foto dan assign ke produk
func (h *UploadHandler) UploadMultipleAndAssignToProduct(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Ambil product_id dari form data
	productIDStr := c.FormValue("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id harus disertakan",
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id harus berupa angka",
		})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal parsing form data",
		})
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File foto tidak ditemukan. Gunakan field 'photos' untuk upload multiple files",
		})
	}

	// Validasi maksimal 5 foto
	if len(files) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Maksimal 5 foto per upload",
		})
	}

	var uploadedFiles []map[string]interface{}
	var photoURLs []string
	var errors []string

	for i, file := range files {
		// Validasi ukuran file (maksimal 5MB)
		if file.Size > 5*1024*1024 {
			errors = append(errors, fmt.Sprintf("File %d: Ukuran terlalu besar (maksimal 5MB)", i+1))
			continue
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
			errors = append(errors, fmt.Sprintf("File %d: Tipe file tidak didukung", i+1))
			continue
		}

		// Generate nama file unik
		timestamp := time.Now().Format("20060102_150405")
		uniqueID := uuid.New().String()[:8]
		fileName := fmt.Sprintf("product_%d_%s_%s_%d%s", uint(userID), timestamp, uniqueID, i+1, fileExt)

		// Path untuk menyimpan file
		uploadPath := filepath.Join("uploads", "products", fileName)

		// Simpan file
		if err := c.SaveFile(file, uploadPath); err != nil {
			errors = append(errors, fmt.Sprintf("File %d: Gagal menyimpan file", i+1))
			continue
		}

		// Tambahkan ke list file yang berhasil diupload
		fileURL := fmt.Sprintf("/uploads/products/%s", fileName)
		photoURLs = append(photoURLs, fileURL)
		uploadedFiles = append(uploadedFiles, map[string]interface{}{
			"filename": fileName,
			"url":      fileURL,
			"size":     file.Size,
		})
	}

	// Assign semua foto ke produk jika ada yang berhasil diupload
	var fotoProduks []interface{}
	if len(photoURLs) > 0 {
		fotoProduks_result, err := h.fotoProdukService.AddMultiplePhotosToProduct(uint(productID), uint(userID), photoURLs)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		
		for _, fp := range fotoProduks_result {
			fotoProduks = append(fotoProduks, fp)
		}
	}

	// Response
	response := fiber.Map{
		"message":        fmt.Sprintf("Berhasil upload %d dari %d file dan ditambahkan ke produk", len(uploadedFiles), len(files)),
		"uploaded_files": uploadedFiles,
		"foto_produks":   fotoProduks,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	statusCode := fiber.StatusCreated
	if len(uploadedFiles) == 0 {
		statusCode = fiber.StatusBadRequest
	} else if len(errors) > 0 {
		statusCode = fiber.StatusPartialContent
	}

	return c.Status(statusCode).JSON(response)
}

// DeleteProductPhoto menghapus foto produk berdasarkan ID foto
func (h *UploadHandler) DeleteProductPhoto(c *fiber.Ctx) error {
	// Ambil userID dari context (dari middleware auth)
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Ambil foto_id dari parameter URL
	fotoIDStr := c.Params("foto_id")
	fotoID, err := strconv.ParseUint(fotoIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID foto tidak valid",
		})
	}

	// Hapus foto menggunakan service
	err = h.fotoProdukService.DeletePhoto(uint(fotoID), uint(userID))
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "tidak memiliki akses") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus foto produk",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Foto produk berhasil dihapus",
	})
}

// GetProductPhotos mengambil semua foto dari produk tertentu
func (h *UploadHandler) GetProductPhotos(c *fiber.Ctx) error {
	// Ambil product_id dari parameter URL
	productIDStr := c.Params("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID produk tidak valid",
		})
	}

	// Ambil foto-foto produk menggunakan service
	photos, err := h.fotoProdukService.GetPhotosByProductID(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil foto produk",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Foto produk berhasil diambil",
		"data":    photos,
	})
}