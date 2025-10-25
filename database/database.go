// file: database/database.go

package database

import (
	"evernos-api2/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3307)/evernos_db1?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
		os.Exit(1) // Keluar dari program dengan status error
	}
	fmt.Println("✅ Connection Opened to Database")
	DB = db
}

func MigrateDB() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Alamat{},
		&models.Toko{},
		&models.Category{},
		&models.Produk{},
		&models.FotoProduk{},
		&models.Trx{},
		&models.DetailTrx{},
		&models.LogProduk{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database!", err)
		os.Exit(1)
	}

	fmt.Println("👍 Database Migration successful")
}
