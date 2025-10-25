package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
	"errors"
	"strconv"
	"strings"
)

type TokoService struct {
	tokoRepo *repositories.TokoRepository
}

func NewTokoService(tokoRepo *repositories.TokoRepository) *TokoService {
	return &TokoService{tokoRepo: tokoRepo}
}

// GetMyToko mengambil toko milik user yang sedang login
func (s *TokoService) GetMyToko(userID uint) (*models.Toko, error) {
	toko, err := s.tokoRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("toko tidak ditemukan")
	}
	return toko, nil
}

// GetTokoByID mengambil toko berdasarkan ID
func (s *TokoService) GetTokoByID(id uint) (*models.Toko, error) {
	toko, err := s.tokoRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("toko tidak ditemukan")
	}
	return toko, nil
}

// UpdateToko memperbarui toko (hanya pemilik yang bisa update)
func (s *TokoService) UpdateToko(id uint, userID uint, updateData map[string]interface{}) (*models.Toko, error) {
	// Cek ownership
	isOwner, err := s.tokoRepo.CheckOwnership(id, userID)
	if err != nil {
		return nil, errors.New("gagal mengecek kepemilikan toko")
	}
	if !isOwner {
		return nil, errors.New("anda tidak memiliki akses untuk mengupdate toko ini")
	}

	// Ambil toko yang akan diupdate
	toko, err := s.tokoRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("toko tidak ditemukan")
	}

	// Validasi dan update field yang diizinkan
	if namaToko, ok := updateData["nama_toko"].(string); ok {
		if err := s.validateNamaToko(namaToko); err != nil {
			return nil, err
		}
		toko.NamaToko = namaToko
	}

	if urlToko, ok := updateData["url_toko"].(string); ok {
		if err := s.validateUrlToko(urlToko); err != nil {
			return nil, err
		}
		toko.UrlToko = urlToko
	}

	// Simpan perubahan
	err = s.tokoRepo.Update(toko)
	if err != nil {
		return nil, errors.New("gagal mengupdate toko")
	}

	return toko, nil
}

// GetAllTokos mengambil semua toko dengan pagination dan filter
func (s *TokoService) GetAllTokos(limitStr, pageStr, namaToko string) ([]models.Toko, map[string]interface{}, error) {
	// Parse limit dan page
	limit := 10 // default
	page := 1   // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Hitung offset
	offset := (page - 1) * limit

	// Ambil data dari repository
	tokos, total, err := s.tokoRepo.GetAllWithPagination(limit, offset, namaToko)
	if err != nil {
		return nil, nil, errors.New("gagal mengambil data toko")
	}

	// Hitung pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := map[string]interface{}{
		"current_page": page,
		"total_pages":  totalPages,
		"total_items":  total,
		"limit":        limit,
		"has_next":     hasNext,
		"has_prev":     hasPrev,
	}

	return tokos, pagination, nil
}

// validateNamaToko memvalidasi nama toko
func (s *TokoService) validateNamaToko(namaToko string) error {
	namaToko = strings.TrimSpace(namaToko)
	if namaToko == "" {
		return errors.New("nama toko tidak boleh kosong")
	}
	if len(namaToko) < 3 {
		return errors.New("nama toko minimal 3 karakter")
	}
	if len(namaToko) > 255 {
		return errors.New("nama toko maksimal 255 karakter")
	}
	return nil
}

// validateUrlToko memvalidasi URL toko
func (s *TokoService) validateUrlToko(urlToko string) error {
	urlToko = strings.TrimSpace(urlToko)
	if urlToko == "" {
		return errors.New("URL toko tidak boleh kosong")
	}
	if len(urlToko) < 3 {
		return errors.New("URL toko minimal 3 karakter")
	}
	if len(urlToko) > 255 {
		return errors.New("URL toko maksimal 255 karakter")
	}
	// Bisa ditambahkan validasi format URL jika diperlukan
	return nil
}

// CreateToko membuat toko baru untuk user
func (s *TokoService) CreateToko(userID uint, tokoData map[string]interface{}) (*models.Toko, error) {
	// Cek apakah user sudah memiliki toko
	existingToko, _ := s.tokoRepo.GetByUserID(userID)
	if existingToko != nil {
		return nil, errors.New("user sudah memiliki toko")
	}

	// Validasi input
	namaToko, ok := tokoData["nama_toko"].(string)
	if !ok {
		return nil, errors.New("nama toko harus diisi")
	}
	if err := s.validateNamaToko(namaToko); err != nil {
		return nil, err
	}

	urlToko, ok := tokoData["url_toko"].(string)
	if !ok {
		return nil, errors.New("URL toko harus diisi")
	}
	if err := s.validateUrlToko(urlToko); err != nil {
		return nil, err
	}

	// Buat toko baru
	toko := &models.Toko{
		IdUser:   userID,
		NamaToko: namaToko,
		UrlToko:  urlToko,
	}

	err := s.tokoRepo.Create(toko)
	if err != nil {
		return nil, errors.New("gagal membuat toko")
	}

	return toko, nil
}

// GetTokoIDByUserID mengambil ID toko berdasarkan user ID
func (s *TokoService) GetTokoIDByUserID(userID uint) (uint, error) {
	toko, err := s.tokoRepo.GetByUserID(userID)
	if err != nil {
		return 0, errors.New("user belum memiliki toko")
	}
	return toko.ID, nil
}

// HasToko mengecek apakah user sudah memiliki toko
func (s *TokoService) HasToko(userID uint) bool {
	toko, err := s.tokoRepo.GetByUserID(userID)
	return err == nil && toko != nil
}