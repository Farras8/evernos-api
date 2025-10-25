package repositories

import (
	"evernos-api2/models"
	"gorm.io/gorm"
)

type TokoRepository struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) *TokoRepository {
	return &TokoRepository{db: db}
}

// GetByUserID mengambil toko berdasarkan user ID (my toko)
func (r *TokoRepository) GetByUserID(userID uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.Where("id_user = ?", userID).First(&toko).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

// GetByID mengambil toko berdasarkan ID
func (r *TokoRepository) GetByID(id uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.First(&toko, id).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

// Create membuat toko baru
func (r *TokoRepository) Create(toko *models.Toko) error {
	return r.db.Create(toko).Error
}

// Update memperbarui toko
func (r *TokoRepository) Update(toko *models.Toko) error {
	return r.db.Save(toko).Error
}

// GetAllWithPagination mengambil semua toko dengan pagination dan filter nama
func (r *TokoRepository) GetAllWithPagination(limit, offset int, namaToko string) ([]models.Toko, int64, error) {
	var tokos []models.Toko
	var total int64

	query := r.db.Model(&models.Toko{})

	// Filter berdasarkan nama toko jika ada (case-insensitive contains)
	if namaToko != "" {
		query = query.Where("LOWER(nama_toko) LIKE LOWER(?)", "%"+namaToko+"%")
	}

	// Hitung total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	err = query.Limit(limit).Offset(offset).Find(&tokos).Error
	if err != nil {
		return nil, 0, err
	}

	return tokos, total, nil
}

// CheckExists mengecek apakah toko dengan ID tertentu ada
func (r *TokoRepository) CheckExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Toko{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// CheckOwnership mengecek apakah toko dengan ID tertentu milik user tertentu
func (r *TokoRepository) CheckOwnership(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Toko{}).Where("id = ? AND id_user = ?", id, userID).Count(&count).Error
	return count > 0, err
}