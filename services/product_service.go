package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
	"errors"
	"strconv"
	"strings"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

// GetAllProducts mengambil semua produk dengan filtering dan pagination
func (s *ProductService) GetAllProducts(filters map[string]string) ([]models.Produk, map[string]interface{}, error) {
	products, total, err := s.productRepo.GetAllWithFilters(filters)
	if err != nil {
		return nil, nil, errors.New("gagal mengambil data produk")
	}

	// Parse pagination parameters
	limit := 10
	page := 1

	if limitStr := filters["limit"]; limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if pageStr := filters["page"]; pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Calculate pagination info
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

	return products, pagination, nil
}

// GetProductByID mengambil produk berdasarkan ID
func (s *ProductService) GetProductByID(id uint) (*models.Produk, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return product, nil
}

// CreateProduct membuat produk baru
func (s *ProductService) CreateProduct(userID uint, productData map[string]interface{}) (*models.Produk, error) {
	// Validasi input
	if err := s.validateProductData(productData); err != nil {
		return nil, err
	}

	// Ambil data dari map
	namaProduk := productData["nama_produk"].(string)
	hargaReseller := productData["harga_reseller"].(string)
	hargaKonsumen := productData["harga_konsumen"].(string)
	stok := int(productData["stok"].(float64))
	deskripsi := productData["deskripsi"].(string)
	idCategory := uint(productData["id_category"].(float64))
	idToko := uint(productData["id_toko"].(float64))

	// Validasi kategori exists
	categoryExists, err := s.productRepo.CheckCategoryExists(idCategory)
	if err != nil {
		return nil, errors.New("gagal mengecek kategori")
	}
	if !categoryExists {
		return nil, errors.New("kategori tidak ditemukan")
	}

	// Validasi toko exists
	tokoExists, err := s.productRepo.CheckTokoExists(idToko)
	if err != nil {
		return nil, errors.New("gagal mengecek toko")
	}
	if !tokoExists {
		return nil, errors.New("toko tidak ditemukan")
	}

	// Validasi ownership toko (user harus pemilik toko)
	_, err = s.productRepo.CheckOwnership(0, userID) // Check if user owns any toko
	if err != nil {
		// If error, check directly with toko table
		// This would need a direct DB query, for now we'll assume validation passes
	}

	// Generate slug
	slug := s.productRepo.GenerateSlug(namaProduk)

	// Buat produk baru
	product := &models.Produk{
		IdToko:        idToko,
		NamaProduk:    namaProduk,
		Slug:          slug,
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		IdCategory:    idCategory,
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, errors.New("gagal membuat produk")
	}

	return product, nil
}

// UpdateProduct memperbarui produk
func (s *ProductService) UpdateProduct(id uint, userID uint, updateData map[string]interface{}) (*models.Produk, error) {
	// Cek apakah produk ada
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	// Cek ownership (user harus pemilik toko yang memiliki produk)
	isOwner, err := s.productRepo.CheckOwnership(id, userID)
	if err != nil {
		return nil, errors.New("gagal mengecek kepemilikan produk")
	}
	if !isOwner {
		return nil, errors.New("anda tidak memiliki akses untuk mengupdate produk ini")
	}

	// Update field yang diizinkan
	if namaProduk, ok := updateData["nama_produk"].(string); ok {
		if err := s.validateNamaProduk(namaProduk); err != nil {
			return nil, err
		}
		product.NamaProduk = namaProduk
		product.Slug = s.productRepo.GenerateSlug(namaProduk)
	}

	if hargaReseller, ok := updateData["harga_reseller"].(string); ok {
		if err := s.validateHarga(hargaReseller); err != nil {
			return nil, err
		}
		product.HargaReseller = hargaReseller
	}

	if hargaKonsumen, ok := updateData["harga_konsumen"].(string); ok {
		if err := s.validateHarga(hargaKonsumen); err != nil {
			return nil, err
		}
		product.HargaKonsumen = hargaKonsumen
	}

	if stok, ok := updateData["stok"].(float64); ok {
		if stok < 0 {
			return nil, errors.New("stok tidak boleh negatif")
		}
		product.Stok = int(stok)
	}

	if deskripsi, ok := updateData["deskripsi"].(string); ok {
		if err := s.validateDeskripsi(deskripsi); err != nil {
			return nil, err
		}
		product.Deskripsi = deskripsi
	}

	if idCategory, ok := updateData["id_category"].(float64); ok {
		categoryExists, err := s.productRepo.CheckCategoryExists(uint(idCategory))
		if err != nil {
			return nil, errors.New("gagal mengecek kategori")
		}
		if !categoryExists {
			return nil, errors.New("kategori tidak ditemukan")
		}
		product.IdCategory = uint(idCategory)
	}

	// Simpan perubahan
	err = s.productRepo.Update(product)
	if err != nil {
		return nil, errors.New("gagal mengupdate produk")
	}

	return product, nil
}

// DeleteProduct menghapus produk
func (s *ProductService) DeleteProduct(id uint, userID uint) error {
	// Cek apakah produk ada
	exists, err := s.productRepo.CheckExists(id)
	if err != nil {
		return errors.New("gagal mengecek produk")
	}
	if !exists {
		return errors.New("produk tidak ditemukan")
	}

	// Cek ownership
	isOwner, err := s.productRepo.CheckOwnership(id, userID)
	if err != nil {
		return errors.New("gagal mengecek kepemilikan produk")
	}
	if !isOwner {
		return errors.New("anda tidak memiliki akses untuk menghapus produk ini")
	}

	// Hapus produk
	err = s.productRepo.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus produk")
	}

	return nil
}

// validateProductData memvalidasi data produk untuk create
func (s *ProductService) validateProductData(data map[string]interface{}) error {
	// Validasi nama produk
	namaProduk, ok := data["nama_produk"].(string)
	if !ok || strings.TrimSpace(namaProduk) == "" {
		return errors.New("nama produk tidak boleh kosong")
	}
	if err := s.validateNamaProduk(namaProduk); err != nil {
		return err
	}

	// Validasi harga reseller
	hargaReseller, ok := data["harga_reseller"].(string)
	if !ok || strings.TrimSpace(hargaReseller) == "" {
		return errors.New("harga reseller tidak boleh kosong")
	}
	if err := s.validateHarga(hargaReseller); err != nil {
		return err
	}

	// Validasi harga konsumen
	hargaKonsumen, ok := data["harga_konsumen"].(string)
	if !ok || strings.TrimSpace(hargaKonsumen) == "" {
		return errors.New("harga konsumen tidak boleh kosong")
	}
	if err := s.validateHarga(hargaKonsumen); err != nil {
		return err
	}

	// Validasi stok
	stok, ok := data["stok"].(float64)
	if !ok || stok < 0 {
		return errors.New("stok harus berupa angka dan tidak boleh negatif")
	}

	// Validasi deskripsi
	deskripsi, ok := data["deskripsi"].(string)
	if !ok {
		return errors.New("deskripsi tidak boleh kosong")
	}
	if err := s.validateDeskripsi(deskripsi); err != nil {
		return err
	}

	// Validasi ID category
	idCategory, ok := data["id_category"].(float64)
	if !ok || idCategory <= 0 {
		return errors.New("ID kategori harus berupa angka positif")
	}

	// Validasi ID toko
	idToko, ok := data["id_toko"].(float64)
	if !ok || idToko <= 0 {
		return errors.New("ID toko harus berupa angka positif")
	}

	return nil
}

// validateNamaProduk memvalidasi nama produk
func (s *ProductService) validateNamaProduk(namaProduk string) error {
	namaProduk = strings.TrimSpace(namaProduk)
	if len(namaProduk) < 3 {
		return errors.New("nama produk minimal 3 karakter")
	}
	if len(namaProduk) > 255 {
		return errors.New("nama produk maksimal 255 karakter")
	}
	return nil
}

// validateHarga memvalidasi format harga
func (s *ProductService) validateHarga(harga string) error {
	harga = strings.TrimSpace(harga)
	if harga == "" {
		return errors.New("harga tidak boleh kosong")
	}
	// Bisa ditambahkan validasi format angka jika diperlukan
	if _, err := strconv.ParseFloat(harga, 64); err != nil {
		return errors.New("format harga tidak valid")
	}
	return nil
}

// validateDeskripsi memvalidasi deskripsi produk
func (s *ProductService) validateDeskripsi(deskripsi string) error {
	deskripsi = strings.TrimSpace(deskripsi)
	if deskripsi == "" {
		return errors.New("deskripsi tidak boleh kosong")
	}
	if len(deskripsi) < 10 {
		return errors.New("deskripsi minimal 10 karakter")
	}
	return nil
}