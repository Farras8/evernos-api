// file: models/entity.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Nama         string    `gorm:"type:varchar(255)"`
	KataSandi    string    `gorm:"type:varchar(255)"`
	NoTelp       string    `gorm:"type:varchar(255);unique"`
	TanggalLahir time.Time `gorm:"type:date"`
	JenisKelamin string    `gorm:"type:varchar(255)"`
	Tentang      string    `gorm:"type:text"`
	Pekerjaan    string    `gorm:"type:varchar(255)"`
	Email        string    `gorm:"type:varchar(255);unique"`
	IdProvinsi   string    `gorm:"type:varchar(255)"`
	IdKota       string    `gorm:"type:varchar(255)"`
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Alamat       []Alamat `gorm:"foreignKey:IdUser"`
	Toko         Toko     `gorm:"foreignKey:IdUser"`
	Trx          []Trx    `gorm:"foreignKey:IdUser"`
}

type Alamat struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	IdUser       uint   `json:"id_user"`
	JudulAlamat  string `gorm:"type:varchar(255)" json:"judul_alamat"`
	NamaPenerima string `gorm:"type:varchar(255)" json:"nama_penerima"`
	NoTelp       string `gorm:"type:varchar(255)" json:"no_telp"`
	DetailAlamat string `gorm:"type:varchar(255)" json:"detail_alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Toko struct {
	ID        uint     `gorm:"primaryKey"`
	IdUser    uint     `gorm:"unique"`
	NamaToko  string   `gorm:"type:varchar(255)"`
	UrlToko   string   `gorm:"type:varchar(255)"`
	Produk    []Produk `gorm:"foreignKey:IdToko"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Produk struct {
	gorm.Model
	IdToko        uint
	NamaProduk    string `gorm:"type:varchar(255)"`
	Slug          string `gorm:"type:varchar(255)"`
	HargaReseller string `gorm:"type:varchar(255)"`
	HargaKonsumen string `gorm:"type:varchar(255)"`
	Stok          int
	Deskripsi     string `gorm:"type:text"`
	IdCategory    uint
	FotoProduk    []FotoProduk `gorm:"foreignKey:IdProduk"`
}

type FotoProduk struct {
	gorm.Model
	IdProduk uint
	Url      string `gorm:"type:varchar(255)"`
}

type Category struct {
	gorm.Model
	NamaCategory string   `gorm:"type:varchar(255)"`
	Produk       []Produk `gorm:"foreignKey:IdCategory"`
}

type Trx struct {
	gorm.Model
	IdUser           uint
	AlamatPengiriman int
	HargaTotal       int
	KodeInvoice      string      `gorm:"type:varchar(255)"`
	MethodBayar      string      `gorm:"type:varchar(255)"`
	DetailTrx        []DetailTrx `gorm:"foreignKey:IdTrx"`
}

type DetailTrx struct {
	gorm.Model
	IdTrx      uint
	IdProduk   uint
	Kuantitas  int
	HargaTotal int
	Produk     Produk `gorm:"foreignKey:IdProduk"`
}

type LogProduk struct {
	gorm.Model
	IdProduk      uint
	IdToko        uint
	IdCategory    uint
	NamaProduk    string `gorm:"type:varchar(255)"`
	Slug          string `gorm:"type:varchar(255)"`
	HargaReseller string `gorm:"type:varchar(255)"`
	HargaKonsumen string `gorm:"type:varchar(255)"`
	Deskripsi     string `gorm:"type:text"`
}

// Response structs for create transaction (without product details)
type DetailTrxCreateResponse struct {
	ID         uint      `json:"ID"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
	DeletedAt  *gorm.DeletedAt `json:"DeletedAt"`
	IdTrx      uint      `json:"IdTrx"`
	IdProduk   uint      `json:"IdProduk"`
	Kuantitas  int       `json:"Kuantitas"`
	HargaTotal int       `json:"HargaTotal"`
}

type TrxCreateResponse struct {
	ID               uint                      `json:"ID"`
	CreatedAt        time.Time                 `json:"CreatedAt"`
	UpdatedAt        time.Time                 `json:"UpdatedAt"`
	DeletedAt        *gorm.DeletedAt           `json:"DeletedAt"`
	IdUser           uint                      `json:"IdUser"`
	AlamatPengiriman int                       `json:"AlamatPengiriman"`
	HargaTotal       int                       `json:"HargaTotal"`
	KodeInvoice      string                    `json:"KodeInvoice"`
	MethodBayar      string                    `json:"MethodBayar"`
	DetailTrx        []DetailTrxCreateResponse `json:"DetailTrx"`
}
