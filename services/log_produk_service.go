package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
)

type LogProdukService struct {
	logProdukRepo *repositories.LogProdukRepository
	productRepo   *repositories.ProductRepository
}

func NewLogProdukService(logProdukRepo *repositories.LogProdukRepository, productRepo *repositories.ProductRepository) *LogProdukService {
	return &LogProdukService{
		logProdukRepo: logProdukRepo,
		productRepo:   productRepo,
	}
}

// LogProductSnapshot mencatat snapshot produk saat transaksi
func (s *LogProdukService) LogProductSnapshot(produkID uint) error {
	// Ambil data produk saat ini
	produk, err := s.productRepo.GetByID(produkID)
	if err != nil {
		return err
	}

	// Buat log snapshot
	logProduk := models.LogProduk{
		IdProduk:      produk.ID,
		IdToko:        produk.IdToko,
		IdCategory:    produk.IdCategory,
		NamaProduk:    produk.NamaProduk,
		Slug:          produk.Slug,
		HargaReseller: produk.HargaReseller,
		HargaKonsumen: produk.HargaKonsumen,
		Deskripsi:     produk.Deskripsi,
	}

	return s.logProdukRepo.Create(&logProduk)
}

// GetLogsByProdukID mengambil semua log untuk produk tertentu
func (s *LogProdukService) GetLogsByProdukID(produkID uint) ([]models.LogProduk, error) {
	return s.logProdukRepo.GetByProdukID(produkID)
}

// GetLogsByTokoID mengambil semua log untuk toko tertentu
func (s *LogProdukService) GetLogsByTokoID(tokoID uint) ([]models.LogProduk, error) {
	return s.logProdukRepo.GetByTokoID(tokoID)
}

// GetAllLogs mengambil semua log produk
func (s *LogProdukService) GetAllLogs() ([]models.LogProduk, error) {
	return s.logProdukRepo.GetAll()
}