package repositories

import (
	"evernos-api2/models"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAllWithFilters mengambil semua produk dengan filtering dan pagination
func (r *ProductRepository) GetAllWithFilters(filters map[string]string) ([]models.Produk, int64, error) {
	var products []models.Produk
	var total int64

	query := r.db.Model(&models.Produk{}).Preload("FotoProduk")

	// Apply filters
	if namaProduk := filters["nama_produk"]; namaProduk != "" {
		query = query.Where("nama_produk LIKE ?", "%"+namaProduk+"%")
	}

	if categoryID := filters["category_id"]; categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			query = query.Where("id_category = ?", uint(id))
		}
	}

	if tokoID := filters["toko_id"]; tokoID != "" {
		if id, err := strconv.ParseUint(tokoID, 10, 32); err == nil {
			query = query.Where("id_toko = ?", uint(id))
		}
	}

	// Price range filters
	if minHarga := filters["min_harga"]; minHarga != "" {
		if price, err := strconv.ParseFloat(minHarga, 64); err == nil {
			// Convert HargaKonsumen string to number for comparison
			query = query.Where("CAST(harga_konsumen AS DECIMAL(10,2)) >= ?", price)
		}
	}

	if maxHarga := filters["max_harga"]; maxHarga != "" {
		if price, err := strconv.ParseFloat(maxHarga, 64); err == nil {
			// Convert HargaKonsumen string to number for comparison
			query = query.Where("CAST(harga_konsumen AS DECIMAL(10,2)) <= ?", price)
		}
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	limit := 10 // default
	page := 1   // default

	if limitStr := filters["limit"]; limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if pageStr := filters["page"]; pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * limit
	err = query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetByID mengambil produk berdasarkan ID
func (r *ProductRepository) GetByID(id uint) (*models.Produk, error) {
	var product models.Produk
	err := r.db.Preload("FotoProduk").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create membuat produk baru
func (r *ProductRepository) Create(product *models.Produk) error {
	return r.db.Create(product).Error
}

// Update memperbarui produk
func (r *ProductRepository) Update(product *models.Produk) error {
	return r.db.Save(product).Error
}

// Delete menghapus produk berdasarkan ID
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Produk{}, id).Error
}

// CheckExists mengecek apakah produk dengan ID tertentu ada
func (r *ProductRepository) CheckExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Produk{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// CheckOwnership mengecek apakah produk dengan ID tertentu milik toko user tertentu
func (r *ProductRepository) CheckOwnership(productID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Table("produks").
		Joins("JOIN tokos ON produks.id_toko = tokos.id").
		Where("produks.id = ? AND tokos.id_user = ?", productID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetByTokoID mengambil produk berdasarkan toko ID
func (r *ProductRepository) GetByTokoID(tokoID uint) ([]models.Produk, error) {
	var products []models.Produk
	err := r.db.Where("id_toko = ?", tokoID).Preload("FotoProduk").Find(&products).Error
	return products, err
}

// GenerateSlug membuat slug dari nama produk
func (r *ProductRepository) GenerateSlug(namaProduk string) string {
	slug := strings.ToLower(namaProduk)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (basic implementation)
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, ",", "")
	slug = strings.ReplaceAll(slug, "!", "")
	slug = strings.ReplaceAll(slug, "?", "")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug
}

// CheckCategoryExists mengecek apakah category dengan ID tertentu ada
func (r *ProductRepository) CheckCategoryExists(categoryID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("id = ?", categoryID).Count(&count).Error
	return count > 0, err
}

// CheckTokoExists mengecek apakah toko dengan ID tertentu ada
func (r *ProductRepository) CheckTokoExists(tokoID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Toko{}).Where("id = ?", tokoID).Count(&count).Error
	return count > 0, err
}