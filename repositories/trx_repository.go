package repositories

import (
	"evernos-api2/models"
	"fmt"
	"gorm.io/gorm"
)

type TrxRepository struct {
	db *gorm.DB
}

func NewTrxRepository(db *gorm.DB) *TrxRepository {
	return &TrxRepository{db: db}
}

// GetByUserID mengambil semua transaksi berdasarkan user ID dengan pagination
func (r *TrxRepository) GetByUserID(userID uint, limit, offset int) ([]models.Trx, int64, error) {
	var trxs []models.Trx
	var total int64

	// Count total records
	err := r.db.Model(&models.Trx{}).Where("id_user = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data with preloaded relations
	err = r.db.Where("id_user = ?", userID).
		Preload("DetailTrx").
		Preload("DetailTrx.Produk").
		Preload("DetailTrx.Produk.FotoProduk").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&trxs).Error

	return trxs, total, err
}

// GetByID mengambil transaksi berdasarkan ID dan user ID (untuk security)
func (r *TrxRepository) GetByID(id uint, userID uint) (*models.Trx, error) {
	var trx models.Trx
	err := r.db.Where("id = ? AND id_user = ?", id, userID).
		Preload("DetailTrx").
		Preload("DetailTrx.Produk").
		Preload("DetailTrx.Produk.FotoProduk").
		First(&trx).Error
	if err != nil {
		return nil, err
	}
	return &trx, nil
}

// Create membuat transaksi baru dengan detail transaksi
func (r *TrxRepository) Create(trx *models.Trx) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create main transaction
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		// Update product stock for each detail
		for _, detail := range trx.DetailTrx {
			var produk models.Produk
			if err := tx.First(&produk, detail.IdProduk).Error; err != nil {
				return err
			}

			// Check if stock is sufficient
			if produk.Stok < detail.Kuantitas {
				return gorm.ErrInvalidData // Will be handled as insufficient stock
			}

			// Update stock
			produk.Stok -= detail.Kuantitas
			if err := tx.Save(&produk).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CheckProductExists mengecek apakah produk dengan ID tertentu ada
func (r *TrxRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Produk{}).Where("id = ?", productID).Count(&count).Error
	return count > 0, err
}

// GetProductByID mengambil produk berdasarkan ID
func (r *TrxRepository) GetProductByID(productID uint) (*models.Produk, error) {
	var produk models.Produk
	err := r.db.First(&produk, productID).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

// CheckAlamatExists mengecek apakah alamat dengan ID tertentu ada untuk user tertentu
func (r *TrxRepository) CheckAlamatExists(alamatID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Alamat{}).Where("id = ? AND id_user = ?", alamatID, userID).Count(&count).Error
	return count > 0, err
}

// GenerateInvoiceCode membuat kode invoice unik
func (r *TrxRepository) GenerateInvoiceCode() string {
	// Simple implementation - in production, use more sophisticated method
	var count int64
	r.db.Model(&models.Trx{}).Count(&count)
	return fmt.Sprintf("INV-%06d", count+1)
}

// CheckExists mengecek apakah transaksi dengan ID tertentu ada untuk user tertentu
func (r *TrxRepository) CheckExists(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Trx{}).Where("id = ? AND id_user = ?", id, userID).Count(&count).Error
	return count > 0, err
}