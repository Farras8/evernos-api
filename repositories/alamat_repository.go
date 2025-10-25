package repositories

import (
	"evernos-api2/models"

	"gorm.io/gorm"
)

type AlamatRepository struct {
	db *gorm.DB
}

func NewAlamatRepository(db *gorm.DB) *AlamatRepository {
	return &AlamatRepository{db: db}
}

// GetByUserID mengambil semua alamat berdasarkan user ID
func (r *AlamatRepository) GetByUserID(userID uint) ([]models.Alamat, error) {
	var alamats []models.Alamat
	err := r.db.Where("id_user = ?", userID).Find(&alamats).Error
	return alamats, err
}

// GetByID mengambil alamat berdasarkan ID dan user ID (untuk security)
func (r *AlamatRepository) GetByID(id uint, userID uint) (*models.Alamat, error) {
	var alamat models.Alamat
	err := r.db.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error
	if err != nil {
		return nil, err
	}
	return &alamat, nil
}

// Create membuat alamat baru
func (r *AlamatRepository) Create(alamat *models.Alamat) error {
	return r.db.Create(alamat).Error
}

// Update memperbarui alamat
func (r *AlamatRepository) Update(alamat *models.Alamat) error {
	return r.db.Save(alamat).Error
}

// Delete menghapus alamat berdasarkan ID dan user ID (untuk security)
func (r *AlamatRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND id_user = ?", id, userID).Delete(&models.Alamat{}).Error
}

// CheckExists mengecek apakah alamat dengan ID tertentu ada untuk user tertentu
func (r *AlamatRepository) CheckExists(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Alamat{}).Where("id = ? AND id_user = ?", id, userID).Count(&count).Error
	return count > 0, err
}
