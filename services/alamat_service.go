package services

import (
	"errors"
	"evernos-api2/models"
	"evernos-api2/repositories"
	"strings"
)

type AlamatService struct {
	alamatRepo *repositories.AlamatRepository
}

func NewAlamatService(alamatRepo *repositories.AlamatRepository) *AlamatService {
	return &AlamatService{alamatRepo: alamatRepo}
}

// GetUserAlamats mengambil semua alamat user
func (s *AlamatService) GetUserAlamats(userID uint) ([]models.Alamat, error) {
	return s.alamatRepo.GetByUserID(userID)
}

// GetAlamatByID mengambil alamat berdasarkan ID untuk user tertentu
func (s *AlamatService) GetAlamatByID(id uint, userID uint) (*models.Alamat, error) {
	alamat, err := s.alamatRepo.GetByID(id, userID)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}
	return alamat, nil
}

// CreateAlamat membuat alamat baru
func (s *AlamatService) CreateAlamat(alamat *models.Alamat) error {
	// Validasi input
	if err := s.validateAlamat(alamat); err != nil {
		return err
	}

	return s.alamatRepo.Create(alamat)
}

// UpdateAlamat memperbarui alamat
func (s *AlamatService) UpdateAlamat(id uint, userID uint, alamat *models.Alamat) error {
	// Cek apakah alamat ada dan milik user
	exists, err := s.alamatRepo.CheckExists(id, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("alamat tidak ditemukan")
	}

	// Ambil data alamat yang sudah ada
	existingAlamat, err := s.alamatRepo.GetByID(id, userID)
	if err != nil {
		return err
	}

	// Update hanya field yang diberikan (partial update)
	if strings.TrimSpace(alamat.JudulAlamat) != "" {
		existingAlamat.JudulAlamat = alamat.JudulAlamat
	}
	if strings.TrimSpace(alamat.NamaPenerima) != "" {
		existingAlamat.NamaPenerima = alamat.NamaPenerima
	}
	if strings.TrimSpace(alamat.NoTelp) != "" {
		existingAlamat.NoTelp = alamat.NoTelp
	}
	if strings.TrimSpace(alamat.DetailAlamat) != "" {
		existingAlamat.DetailAlamat = alamat.DetailAlamat
	}

	// Validasi data yang sudah diupdate
	if err := s.validateAlamat(existingAlamat); err != nil {
		return err
	}

	return s.alamatRepo.Update(existingAlamat)
}

// DeleteAlamat menghapus alamat
func (s *AlamatService) DeleteAlamat(id uint, userID uint) error {
	// Cek apakah alamat ada dan milik user
	exists, err := s.alamatRepo.CheckExists(id, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("alamat tidak ditemukan")
	}

	return s.alamatRepo.Delete(id, userID)
}

// validateAlamat melakukan validasi input alamat
func (s *AlamatService) validateAlamat(alamat *models.Alamat) error {
	if strings.TrimSpace(alamat.JudulAlamat) == "" {
		return errors.New("judul alamat tidak boleh kosong")
	}
	if strings.TrimSpace(alamat.NamaPenerima) == "" {
		return errors.New("nama penerima tidak boleh kosong")
	}
	if strings.TrimSpace(alamat.NoTelp) == "" {
		return errors.New("nomor telepon tidak boleh kosong")
	}
	if strings.TrimSpace(alamat.DetailAlamat) == "" {
		return errors.New("detail alamat tidak boleh kosong")
	}
	return nil
}
