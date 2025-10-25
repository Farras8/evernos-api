package repositories

import (
	"evernos-api2/models"
	"gorm.io/gorm"
)

type FotoProdukRepository struct {
	db *gorm.DB
}

func NewFotoProdukRepository(db *gorm.DB) *FotoProdukRepository {
	return &FotoProdukRepository{db: db}
}

// Create menambahkan foto produk baru
func (r *FotoProdukRepository) Create(fotoProduk *models.FotoProduk) error {
	return r.db.Create(fotoProduk).Error
}

// CreateMultiple menambahkan multiple foto produk
func (r *FotoProdukRepository) CreateMultiple(fotoProduks []models.FotoProduk) error {
	return r.db.Create(&fotoProduks).Error
}

// GetByProductID mengambil semua foto berdasarkan product ID
func (r *FotoProdukRepository) GetByProductID(productID uint) ([]models.FotoProduk, error) {
	var fotoProduks []models.FotoProduk
	err := r.db.Where("id_produk = ?", productID).Find(&fotoProduks).Error
	return fotoProduks, err
}

// DeleteByID menghapus foto berdasarkan ID
func (r *FotoProdukRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.FotoProduk{}, id).Error
}

// DeleteByProductID menghapus semua foto berdasarkan product ID
func (r *FotoProdukRepository) DeleteByProductID(productID uint) error {
	return r.db.Where("id_produk = ?", productID).Delete(&models.FotoProduk{}).Error
}

// CheckExists mengecek apakah foto dengan ID tertentu ada
func (r *FotoProdukRepository) CheckExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.FotoProduk{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// CheckOwnership mengecek apakah foto dengan ID tertentu milik produk user tertentu
func (r *FotoProdukRepository) CheckOwnership(fotoID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Table("foto_produks").
		Joins("JOIN produks ON foto_produks.id_produk = produks.id").
		Joins("JOIN tokos ON produks.id_toko = tokos.id").
		Where("foto_produks.id = ? AND tokos.id_user = ?", fotoID, userID).
		Count(&count).Error
	return count > 0, err
}