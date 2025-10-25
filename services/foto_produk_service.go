package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
	"errors"
)

type FotoProdukService struct {
	fotoProdukRepo *repositories.FotoProdukRepository
	productRepo    *repositories.ProductRepository
}

func NewFotoProdukService(fotoProdukRepo *repositories.FotoProdukRepository, productRepo *repositories.ProductRepository) *FotoProdukService {
	return &FotoProdukService{
		fotoProdukRepo: fotoProdukRepo,
		productRepo:    productRepo,
	}
}

// AddPhotoToProduct menambahkan foto ke produk
func (s *FotoProdukService) AddPhotoToProduct(productID uint, userID uint, photoURL string) (*models.FotoProduk, error) {
	// Cek apakah produk ada
	productExists, err := s.productRepo.CheckExists(productID)
	if err != nil {
		return nil, errors.New("gagal mengecek produk")
	}
	if !productExists {
		return nil, errors.New("produk tidak ditemukan")
	}

	// Cek ownership (user harus pemilik toko yang memiliki produk)
	isOwner, err := s.productRepo.CheckOwnership(productID, userID)
	if err != nil {
		return nil, errors.New("gagal mengecek kepemilikan produk")
	}
	if !isOwner {
		return nil, errors.New("anda tidak memiliki akses untuk menambahkan foto ke produk ini")
	}

	// Buat foto produk baru
	fotoProduk := &models.FotoProduk{
		IdProduk: productID,
		Url:      photoURL,
	}

	err = s.fotoProdukRepo.Create(fotoProduk)
	if err != nil {
		return nil, errors.New("gagal menyimpan foto produk")
	}

	return fotoProduk, nil
}

// AddMultiplePhotosToProduct menambahkan multiple foto ke produk
func (s *FotoProdukService) AddMultiplePhotosToProduct(productID uint, userID uint, photoURLs []string) ([]models.FotoProduk, error) {
	// Cek apakah produk ada
	productExists, err := s.productRepo.CheckExists(productID)
	if err != nil {
		return nil, errors.New("gagal mengecek produk")
	}
	if !productExists {
		return nil, errors.New("produk tidak ditemukan")
	}

	// Cek ownership (user harus pemilik toko yang memiliki produk)
	isOwner, err := s.productRepo.CheckOwnership(productID, userID)
	if err != nil {
		return nil, errors.New("gagal mengecek kepemilikan produk")
	}
	if !isOwner {
		return nil, errors.New("anda tidak memiliki akses untuk menambahkan foto ke produk ini")
	}

	// Buat array foto produk
	var fotoProduks []models.FotoProduk
	for _, url := range photoURLs {
		fotoProduks = append(fotoProduks, models.FotoProduk{
			IdProduk: productID,
			Url:      url,
		})
	}

	err = s.fotoProdukRepo.CreateMultiple(fotoProduks)
	if err != nil {
		return nil, errors.New("gagal menyimpan foto produk")
	}

	return fotoProduks, nil
}

// GetPhotosByProductID mengambil semua foto berdasarkan product ID
func (s *FotoProdukService) GetPhotosByProductID(productID uint) ([]models.FotoProduk, error) {
	return s.fotoProdukRepo.GetByProductID(productID)
}

// DeletePhoto menghapus foto produk
func (s *FotoProdukService) DeletePhoto(fotoID uint, userID uint) error {
	// Cek apakah foto ada
	fotoExists, err := s.fotoProdukRepo.CheckExists(fotoID)
	if err != nil {
		return errors.New("gagal mengecek foto")
	}
	if !fotoExists {
		return errors.New("foto tidak ditemukan")
	}

	// Cek ownership (user harus pemilik toko yang memiliki produk)
	isOwner, err := s.fotoProdukRepo.CheckOwnership(fotoID, userID)
	if err != nil {
		return errors.New("gagal mengecek kepemilikan foto")
	}
	if !isOwner {
		return errors.New("anda tidak memiliki akses untuk menghapus foto ini")
	}

	err = s.fotoProdukRepo.DeleteByID(fotoID)
	if err != nil {
		return errors.New("gagal menghapus foto")
	}

	return nil
}