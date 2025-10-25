package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
	"errors"
	"strconv"
	"strings"
	"gorm.io/gorm"
)

type TrxService struct {
	trxRepo        *repositories.TrxRepository
	logProdukService *LogProdukService
}

func NewTrxService(trxRepo *repositories.TrxRepository, logProdukService *LogProdukService) *TrxService {
	return &TrxService{
		trxRepo:        trxRepo,
		logProdukService: logProdukService,
	}
}

// GetAllTrx mengambil semua transaksi user dengan pagination
func (s *TrxService) GetAllTrx(userID uint, limitStr, pageStr string) ([]models.Trx, map[string]interface{}, error) {
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
	trxs, total, err := s.trxRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, nil, errors.New("gagal mengambil data transaksi")
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

	return trxs, pagination, nil
}

// GetTrxByID mengambil transaksi berdasarkan ID
func (s *TrxService) GetTrxByID(id uint, userID uint) (*models.Trx, error) {
	trx, err := s.trxRepo.GetByID(id, userID)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}
	return trx, nil
}

// CreateTrx membuat transaksi baru
func (s *TrxService) CreateTrx(userID uint, trxData map[string]interface{}) (*models.Trx, error) {
	// Validasi data
	if err := s.validateTrxData(trxData); err != nil {
		return nil, err
	}

	// Extract data
	methodBayar := trxData["method_bayar"].(string)
	alamatKirim := int(trxData["alamat_kirim"].(float64))
	detailTrxData := trxData["detail_trx"].([]interface{})

	// Validasi alamat pengiriman
	alamatExists, err := s.trxRepo.CheckAlamatExists(uint(alamatKirim), userID)
	if err != nil {
		return nil, errors.New("gagal mengecek alamat")
	}
	if !alamatExists {
		return nil, errors.New("alamat pengiriman tidak ditemukan")
	}

	// Validasi dan hitung total harga
	var detailTrxs []models.DetailTrx
	var totalHarga int

	for _, detail := range detailTrxData {
		detailMap := detail.(map[string]interface{})
		productID := uint(detailMap["product_id"].(float64))
		kuantitas := int(detailMap["kuantitas"].(float64))

		// Validasi produk exists
		productExists, err := s.trxRepo.CheckProductExists(productID)
		if err != nil {
			return nil, errors.New("gagal mengecek produk")
		}
		if !productExists {
			return nil, errors.New("produk dengan ID " + strconv.Itoa(int(productID)) + " tidak ditemukan")
		}

		// Ambil data produk untuk hitung harga
		produk, err := s.trxRepo.GetProductByID(productID)
		if err != nil {
			return nil, errors.New("gagal mengambil data produk")
		}

		// Validasi stok
		if produk.Stok < kuantitas {
			return nil, errors.New("stok produk " + produk.NamaProduk + " tidak mencukupi")
		}

		// Hitung harga (menggunakan harga konsumen)
		hargaKonsumen, err := strconv.Atoi(produk.HargaKonsumen)
		if err != nil {
			return nil, errors.New("harga produk tidak valid")
		}

		hargaDetail := hargaKonsumen * kuantitas
		totalHarga += hargaDetail

		// Tambahkan ke detail transaksi
		detailTrxs = append(detailTrxs, models.DetailTrx{
			IdProduk:   productID,
			Kuantitas:  kuantitas,
			HargaTotal: hargaDetail,
		})
	}

	// Generate kode invoice
	kodeInvoice := s.trxRepo.GenerateInvoiceCode()

	// Buat transaksi
	trx := &models.Trx{
		IdUser:           userID,
		AlamatPengiriman: alamatKirim,
		HargaTotal:       totalHarga,
		KodeInvoice:      kodeInvoice,
		MethodBayar:      methodBayar,
		DetailTrx:        detailTrxs,
	}

	err = s.trxRepo.Create(trx)
	if err != nil {
		if err == gorm.ErrInvalidData {
			return nil, errors.New("stok produk tidak mencukupi")
		}
		return nil, errors.New("gagal membuat transaksi")
	}

	// Log produk snapshot untuk setiap produk dalam transaksi
	for _, detail := range detailTrxs {
		if logErr := s.logProdukService.LogProductSnapshot(detail.IdProduk); logErr != nil {
			// Log error tapi jangan gagalkan transaksi
			// Bisa ditambahkan logging ke file atau monitoring system
			continue
		}
	}

	return trx, nil
}

// CreateTrxWithResponse membuat transaksi baru dan mengembalikan response tanpa detail produk
func (s *TrxService) CreateTrxWithResponse(userID uint, trxData map[string]interface{}) (*models.TrxCreateResponse, error) {
	trx, err := s.CreateTrx(userID, trxData)
	if err != nil {
		return nil, err
	}
	
	return s.convertToCreateResponse(trx), nil
}

// convertToCreateResponse mengkonversi Trx ke TrxCreateResponse (tanpa detail produk)
func (s *TrxService) convertToCreateResponse(trx *models.Trx) *models.TrxCreateResponse {
	var detailTrxResponses []models.DetailTrxCreateResponse
	
	for _, detail := range trx.DetailTrx {
		detailTrxResponses = append(detailTrxResponses, models.DetailTrxCreateResponse{
			ID:         detail.ID,
			CreatedAt:  detail.CreatedAt,
			UpdatedAt:  detail.UpdatedAt,
			IdTrx:      detail.IdTrx,
			IdProduk:   detail.IdProduk,
			Kuantitas:  detail.Kuantitas,
			HargaTotal: detail.HargaTotal,
		})
	}
	
	return &models.TrxCreateResponse{
		ID:               trx.ID,
		CreatedAt:        trx.CreatedAt,
		UpdatedAt:        trx.UpdatedAt,
		IdUser:           trx.IdUser,
		AlamatPengiriman: trx.AlamatPengiriman,
		HargaTotal:       trx.HargaTotal,
		KodeInvoice:      trx.KodeInvoice,
		MethodBayar:      trx.MethodBayar,
		DetailTrx:        detailTrxResponses,
	}
}

// validateTrxData memvalidasi data transaksi
func (s *TrxService) validateTrxData(data map[string]interface{}) error {
	// Validasi method_bayar
	methodBayar, ok := data["method_bayar"].(string)
	if !ok || strings.TrimSpace(methodBayar) == "" {
		return errors.New("method bayar tidak boleh kosong")
	}

	// Validasi alamat_kirim
	alamatKirim, ok := data["alamat_kirim"].(float64)
	if !ok || alamatKirim <= 0 {
		return errors.New("alamat kirim tidak valid")
	}

	// Validasi detail_trx
	detailTrx, ok := data["detail_trx"].([]interface{})
	if !ok || len(detailTrx) == 0 {
		return errors.New("detail transaksi tidak boleh kosong")
	}

	// Validasi setiap detail transaksi
	for i, detail := range detailTrx {
		detailMap, ok := detail.(map[string]interface{})
		if !ok {
			return errors.New("format detail transaksi tidak valid")
		}

		// Validasi product_id
		productID, ok := detailMap["product_id"].(float64)
		if !ok || productID <= 0 {
			return errors.New("product_id pada detail ke-" + strconv.Itoa(i+1) + " tidak valid")
		}

		// Validasi kuantitas
		kuantitas, ok := detailMap["kuantitas"].(float64)
		if !ok || kuantitas <= 0 {
			return errors.New("kuantitas pada detail ke-" + strconv.Itoa(i+1) + " tidak valid")
		}
	}

	return nil
}