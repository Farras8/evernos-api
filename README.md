# Evernos API Backend

Evernos API adalah backend service untuk aplikasi e-commerce yang menyediakan berbagai endpoint untuk mengelola produk, kategori, toko, transaksi, dan fitur-fitur lainnya.

## 🚀 Fitur Utama

- **Autentikasi & Autorisasi**: Login, register, dan middleware JWT
- **Manajemen Produk**: CRUD produk dengan upload foto
- **Manajemen Kategori**: CRUD kategori produk
- **Manajemen Toko**: CRUD toko dan profil toko
- **Manajemen Alamat**: CRUD alamat user
- **Sistem Transaksi**: Pembuatan dan pengelolaan transaksi
- **Regional Data**: Data provinsi dan kota Indonesia
- **Upload File**: Upload foto produk dengan validasi

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Fiber v2
- **Database**: MySQL
- **ORM**: GORM
- **Authentication**: JWT
- **File Upload**: Multipart form handling
- **Validation**: Built-in validation

## 📁 Struktur Proyek

```
evernos-api2/
├── database/           # Konfigurasi database dan seeding
├── handlers/           # HTTP handlers untuk setiap endpoint
├── middleware/         # Middleware untuk autentikasi
├── models/            # Entity models dan structs
├── repositories/      # Data access layer
├── routes/            # Route definitions
├── services/          # Business logic layer
├── uploads/           # Direktori untuk file upload
├── main.go            # Entry point aplikasi
├── go.mod             # Go modules
└── .env               # Environment variables
```

## 🔧 Instalasi & Setup

### Prerequisites
- Go 1.19 atau lebih baru
- MySQL/PostgreSQL database
- Git

### Langkah Instalasi

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd evernos-api2
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment variables**
   
   Buat file `.env` di root directory:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=evernos_db
   JWT_SECRET=your_jwt_secret_key
   PORT=3001
   ```

4. **Setup database**
   ```bash
   # Buat database baru
   mysql -u root -p
   CREATE DATABASE evernos_db;
   ```

5. **Run aplikasi**
   ```bash
   go run main.go
   ```

   Server akan berjalan di `http://localhost:3001`

## 📚 API Documentation

Dokumentasi lengkap API tersedia di Postman:

**🔗 [Evernos API Documentation](https://documenter.getpostman.com/view/36349178/2sB3QRp7fB)**

### Endpoint Utama

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/auth/login` | Login user |
| POST | `/auth/register` | Register user baru |
| GET | `/category` | Get semua kategori |
| POST | `/category` | Buat kategori baru (Admin) |
| PUT | `/category/:id` | Update kategori (Admin) |
| DELETE | `/category/:id` | Hapus kategori (Admin) |
| GET | `/product` | Get semua produk |
| POST | `/product` | Buat produk baru |
| PUT | `/product/:id` | Update produk |
| DELETE | `/product/:id` | Hapus produk |
| GET | `/toko` | Get semua toko |
| POST | `/toko` | Buat toko baru |
| PUT | `/toko/:id_toko` | Update toko |
| POST | `/trx` | Buat transaksi baru |
| GET | `/trx` | Get riwayat transaksi |

## 🔐 Autentikasi

API menggunakan JWT (JSON Web Token) untuk autentikasi. Setelah login berhasil, sertakan token di header:

```
Authorization: Bearer <your_jwt_token>
```

### Role-based Access
- **User**: Akses ke produk, toko, alamat, transaksi
- **Admin**: Akses penuh termasuk manajemen kategori

## 📝 Testing

Untuk testing API, gunakan file testing guide yang tersedia:
- `API_TESTING_GUIDE.md` - Panduan lengkap testing semua endpoint
- `product_test_data.json` - Data testing untuk produk
- `category_test.json` - Data testing untuk kategori

### Contoh Testing dengan cURL

```bash
# Login
curl -X POST http://localhost:3001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get semua produk
curl -X GET http://localhost:3001/product \
  -H "Authorization: Bearer <your_token>"

# Buat produk baru
curl -X POST http://localhost:3001/product \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"nama_produk":"Produk Test","harga":50000,"deskripsi":"Deskripsi produk","id_category":"1"}'
```

## 🗂️ Database Schema

### Tabel Utama
- `users` - Data user dan autentikasi
- `categories` - Kategori produk
- `products` - Data produk
- `tokos` - Data toko
- `alamats` - Alamat user
- `trxs` - Transaksi
- `detail_trxs` - Detail item transaksi
- `foto_produks` - Foto produk
- `log_produks` - Log perubahan produk

## 🚀 Deployment

### Development
```bash
go run main.go
```

### Production
```bash
# Build aplikasi
go build -o evernos-api main.go

# Run binary
./evernos-api
```

### Docker (Opsional)
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 🤝 Contributing

1. Fork repository
2. Buat feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📞 Contact

- **Developer**: Evernos Team
- **Email**: support@evernos.com
- **API Documentation**: [Postman Collection](https://documenter.getpostman.com/view/36349178/2sB3QRp7fB)

---

**Happy Coding! 🎉**