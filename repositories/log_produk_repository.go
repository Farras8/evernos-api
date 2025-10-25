package repositories

import (
	"evernos-api2/models"
	"gorm.io/gorm"
)

type LogProdukRepository struct {
	db *gorm.DB
}

func NewLogProdukRepository(db *gorm.DB) *LogProdukRepository {
	return &LogProdukRepository{db: db}
}

func (r *LogProdukRepository) Create(logProduk *models.LogProduk) error {
	return r.db.Create(logProduk).Error
}

func (r *LogProdukRepository) GetByProdukID(produkID uint) ([]models.LogProduk, error) {
	var logs []models.LogProduk
	err := r.db.Where("id_produk = ?", produkID).Find(&logs).Error
	return logs, err
}

func (r *LogProdukRepository) GetByTokoID(tokoID uint) ([]models.LogProduk, error) {
	var logs []models.LogProduk
	err := r.db.Where("id_toko = ?", tokoID).Find(&logs).Error
	return logs, err
}

func (r *LogProdukRepository) GetAll() ([]models.LogProduk, error) {
	var logs []models.LogProduk
	err := r.db.Find(&logs).Error
	return logs, err
}