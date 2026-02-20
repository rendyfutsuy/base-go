# Base Template Go - REST API Backend

Template aplikasi backend REST API yang dibuat menggunakan Go (Golang), Echo Framework, dan GORM. Base template ini menyediakan struktur aplikasi yang siap digunakan untuk pengembangan aplikasi web modern dengan fitur autentikasi, manajemen user, dan manajemen role yang lengkap.

## 📋 Daftar Isi

- [Fitur Utama](#fitur-utama)
- [Persyaratan](#persyaratan)
- [Instalasi](#instalasi)
- [Konfigurasi](#konfigurasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Database Migration & Seeding](#database-migration--seeding)
- [Struktur Proyek](#struktur-proyek)
- [API Documentation](#api-documentation)
- [Docker](#docker)
- [Development Guidelines](#development-guidelines)

## 🚀 Fitur Utama

### 1. Autentikasi & Authorisasi
- ✅ Login/Logout dengan JWT Token
- ✅ Refresh Token
- ✅ Reset Password via Email
- ✅ Update Password dengan validasi password history
- ✅ Password expiration & attempt counter
- ✅ JWT Token Management
- ✅ Profile Management

### 2. User Management
- ✅ CRUD User (Create, Read, Update, Delete)
- ✅ Soft Delete User
- ✅ Pagination & Filtering
- ✅ Bulk Import User dari Excel
- ✅ Download Template Excel untuk Import
- ✅ Validasi duplikasi (Email, Username, NIK)
- ✅ Block/Unblock User
- ✅ Activate/Deactivate User
- ✅ Password Management

### 3. Role Management
- ✅ CRUD Role (Create, Read, Update, Delete)
- ✅ Permission Management
- ✅ Permission Group Management
- ✅ Assign Permission ke Role
- ✅ Assign Permission ke User

### 4. Fitur Pendukung
- ✅ Swagger Documentation (hanya di development)
- ✅ Request Validation dengan Custom Validator
- ✅ Rate Limiting (100 requests/second)
- ✅ CORS Configuration
- ✅ Logging dengan Zap Logger
- ✅ NewRelic Integration (optional)
- ✅ Background Jobs dengan Asynq
- ✅ Redis Support
- ✅ Email Service
- ✅ File Upload Support

## 📦 Persyaratan

Sebelum memulai, pastikan Anda telah menginstall:

- **Go** 1.24.0 atau lebih tinggi
- **PostgreSQL** 12.0 atau lebih tinggi
- **Redis** (optional, untuk background jobs)
- **Migrate CLI** (untuk database migration)
- **Task** (optional, untuk task runner)
- **Swag** (untuk generate Swagger documentation)

### Install Dependencies

```bash
# Install Migrate CLI
brew install migrate  # macOS
# atau download dari: https://github.com/golang-migrate/migrate

# Install Task (optional)
brew install go-task/tap/go-task  # macOS

# Install Swag
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🔧 Instalasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd base-go
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Database

Buat database PostgreSQL:

```bash
createdb base-local
# atau
psql -U postgres -c "CREATE DATABASE base_local;"
```

### 4. Setup Config

Copy file konfigurasi contoh dan sesuaikan dengan kebutuhan:

```bash
cp .env.example .env
```

Edit `.env` sesuai dengan environment Anda:

```
APP_NAME=base v2.0
APP_ENV=development
APP_PORT=9090
QUEUE_PORT=9091
CONTEXT__TIMEOUT=30

JWT_KEY=123456
JWT_REFRESH_KEY=123456

NEWRELIC__ENABLE_NEW_RELIC_LOGGING=false
NEWRELIC__APP_NAME=Go Base Project
NEWRELIC__LICENSE=
NEWRELIC__IDENTIFIER=local

DATABASE__HOST=localhost
DATABASE__PORT=5432
DATABASE__USER=postgres
DATABASE__PASSWORD=postgres
DATABASE__DB_NAME=go_base_project
DATABASE__SSLMODE=disable
DATABASE__TOKEN_STORAGE=local
....
```

## 🔄 Database Migration & Seeding

### Migration

Jalankan migration untuk membuat schema database:

```bash
# Menggunakan migrate CLI
migrate -path ./database/migrations \
  -database "postgresql://USERNAME:PASSWORD@localhost:5432/DATABASE_NAME?sslmode=disable" \
  up

# Contoh:
migrate -path ./database/migrations \
  -database "postgresql://postgres:postgres@localhost:5432/base_local?sslmode=disable" \
  up

# Rollback migration
migrate -path ./database/migrations \
  -database "postgresql://USERNAME:PASSWORD@localhost:5432/DATABASE_NAME?sslmode=disable" \
  down
```

### Seeding

Jalankan seeder untuk mengisi data awal:

```bash
migrate -path ./database/seeders \
  -database "postgresql://USERNAME:PASSWORD@localhost:5432/DATABASE_NAME?sslmode=disable&x-migrations-table=seeder_migrations" \
  up

# Contoh:
migrate -path ./database/seeders \
  -database "postgresql://postgres:postgres@localhost:5432/base_local?sslmode=disable&x-migrations-table=seeder_migrations" \
  up
```

### Menggunakan Taskfile (Optional)

Jika menggunakan Taskfile, sesuaikan konfigurasi di `Taskfile.yml`:

```bash
# Jalankan migration
task migrate-up

# Rollback migration
task migrate-down

# Jalankan seeder
task seed-up

# Rollback seeder
task seed-down
```

## ▶️ Menjalankan Aplikasi

### Development Mode

```bash
# Generate Swagger documentation
swag init -g router/router.go

# Jalankan aplikasi
go run main.go
```

Aplikasi akan berjalan di `http://localhost:9090` (sesuai konfigurasi `app_port` di `.env`).

### Build & Run

```bash
# Build aplikasi
go build -o base-go main.go

# Jalankan aplikasi
./base-go
```

### Menggunakan Air (Hot Reload)

Untuk development dengan hot reload, install Air terlebih dahulu:

```bash
go install github.com/cosmtrek/air@latest
```

Jalankan dengan:

```bash
air
```

## 📚 API Documentation

### Swagger UI

Akses Swagger documentation di development environment:

```
http://localhost:9090/swagger/index.html
```

**Catatan:** Swagger hanya tersedia jika `app_env` di `.env` diset ke `development`.

### Regenerate Swagger Documentation

Setelah menambahkan atau mengubah dokumentasi API:

```bash
swag init -g router/router.go
```

## 🐳 Docker

### Build Docker Image

```bash
docker build -t base-go:latest .
```

### Run dengan Docker Compose

```bash
docker-compose up -d
```

Aplikasi akan berjalan di `http://localhost:9000`.

## 📁 Struktur Proyek

```
base-go/
├── constants/          # Konstanta aplikasi
├── database/           # Database setup, migrations, seeders
│   ├── migrations/     # Database migrations
│   └── seeders/        # Database seeders
├── docs/               # Swagger documentation (auto-generated)
├── helpers/            # Helper functions dan middleware
│   ├── middleware/     # Custom middleware
│   ├── request/        # Request helpers
│   └── response/       # Response helpers
├── models/             # Database models
├── modules/            # Business logic modules
│   ├── auth/           # Authentication module
│   ├── user_management/# User management module
│   ├── role_management/# Role management module
│   └── homepage/       # Homepage module
├── router/             # Router configuration
├── utils/              # Utility functions
├── worker/             # Background jobs
├── public/             # Public files
├── main.go             # Application entry point
├── config.json         # Application configuration
└── go.mod              # Go module dependencies
└── .env                # Env file
```

### Arsitektur Clean Architecture

Proyek ini mengikuti prinsip Clean Architecture dengan pemisahan layer:

```
modules/{module_name}/
├── delivery/           # Presentation layer
│   └── http/           # HTTP handlers
├── usecase/            # Business logic layer
├── repository/         # Data access layer
├── dto/                # Data Transfer Objects
├── repository.go       # Repository interface
└── usecase.go          # Usecase interface
```

## 🔐 Authentication

### Login

```bash
curl -X POST http://localhost:9090/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Menggunakan Token

Setelah login, gunakan token di header Authorization:

```bash
curl -X GET http://localhost:9090/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 👤 User Management

### Create User

```bash
curl -X POST http://localhost:9090/v1/user-management/user \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "username": "newuser",
    "full_name": "New User",
    "nik": "1234567890123456",
    "role_id": "uuid-role-id",
    "password": "password123",
    "is_active": true
  }'
```

### Import Users dari Excel

1. Download template Excel:
```bash
curl -X GET http://localhost:9090/v1/user-management/user/import/template \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  --output user_import_template.xlsx
```

2. Isi template dengan data user

3. Upload file Excel:
```bash
curl -X POST http://localhost:9090/v1/user-management/user/import \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "file=@user_import_template.xlsx"
```

Response jika ada error:
```json
{
  "data": {
    "total_rows": 5,
    "success_count": 3,
    "failed_count": 2,
    "results": [
      {
        "row": 2,
        "username": "user1",
        "status": "failed",
        "error_message": "Email sudah terdaftar di database; Username sudah terdaftar di database"
      },
      {
        "row": 3,
        "username": "user2",
        "status": "success"
      }
    ]
  }
}
```

## 🛠️ Development Guidelines

### Menambahkan Module Baru

1. **Buat struktur folder:**
```
modules/new_module/
├── delivery/
│   └── http/
├── usecase/
├── repository/
├── dto/
├── repository.go
└── usecase.go
```

2. **Define Interface di `repository.go` dan `usecase.go`**

3. **Implement repository di `repository/new_module_repository.go`**

4. **Implement usecase di `usecase/new_module_usecase.go`**

5. **Create HTTP handlers di `delivery/http/new_module_handler.go`**

6. **Register routes di `router/router.go`**

### Menambahkan Migration

1. Buat file migration:
```bash
migrate create -ext sql -dir database/migrations -seq create_new_table
```

2. Edit file `.up.sql` dan `.down.sql`

3. Jalankan migration:
```bash
migrate -path ./database/migrations \
  -database "postgresql://..." \
  up
```

### Menambahkan Model

1. Buat model di `models/new_model.go` dengan GORM tags
2. Implement `TableName()` method jika diperlukan
3. Gunakan model di repository

### Testing

Jalankan unit test:

```bash
go test ./...
```

Jalankan test dengan coverage:

```bash
go test -cover ./...
```

## 🔧 Konfigurasi

### Environment Variables

Aplikasi menggunakan file `.env` untuk konfigurasi. Untuk production, disarankan menggunakan environment variables atau service seperti Vault.

### Konfigurasi yang Tersedia

- `app_name`: Nama aplikasi
- `app_env`: Environment (development, staging, production)
- `app_port`: Port aplikasi
- `database`: Konfigurasi database PostgreSQL
- `redis`: Konfigurasi Redis
- `jwt_key`: Secret key untuk JWT
- `email`: Konfigurasi SMTP
- `auth.access_token_ttl_seconds`: TTL untuk access token (dalam detik)

## 📝 Best Practices

1. **Always validate input** - Gunakan validator untuk memastikan data yang masuk valid
2. **Use constants** - Simpan string magic ke dalam constants
3. **Error handling** - Always handle error dengan proper error messages
4. **Logging** - Gunakan structured logging untuk debugging
5. **Transactions** - Gunakan transaction untuk operasi database yang kompleks
6. **Batch operations** - Untuk bulk operations, gunakan batch processing
7. **Security** - Jangan hardcode credentials, gunakan config atau environment variables

## 🐛 Troubleshooting

### Database Connection Error

Pastikan:
- PostgreSQL sudah berjalan
- Credentials di `.env` benar
- Database sudah dibuat
- Migration sudah dijalankan

### Port Already in Use

Ubah `app_port` di `.env` atau kill process yang menggunakan port tersebut.

### Swagger Not Available

Pastikan `app_env` di `.env` diset ke `development`.

## 📄 License

Rendy Anggara © since 2025

## 👥 Contributors

<table>

<tr>

<td  align="center"><a  href="https://github.com/rendyfutsuy"><img  src="https://avatars.githubusercontent.com/u/22336340?s=96&v=4"  style="border-radius:50%;" width="100px;"  alt=""/><br  /><sub><b>Rendy Anggara</b></sub></a><br  /><a  href="#"  title="Owner">⚜</a><a  href="#"  title="Code">💻</a></td>
</tr>

</table>

## 📞 Support

Untuk pertanyaan atau masalah, silakan buat issue di repository ini.

---

**Happy Coding! 🚀**

