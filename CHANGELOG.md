# Changelog

## [2.0.0] - 2026-02-27

### Added - Priority High Features

#### 🔐 Authentication & Authorization
- JWT-based authentication untuk REST API
- gRPC interceptor untuk autentikasi gRPC requests
- Endpoint `/auth/register` untuk registrasi user baru
- Endpoint `/auth/login` untuk autentikasi
- Password hashing dengan bcrypt
- Protected endpoints yang memerlukan JWT token
- Middleware `JWTAuth` untuk validasi token di REST API
- Middleware `GRPCAuthInterceptor` untuk validasi token di gRPC

#### ⚙️ Environment-based Configuration
- Config management dengan environment variables
- File `.env.example` sebagai template
- Support untuk multiple environments (development/production)
- Konfigurasi untuk:
  - HTTP Port
  - gRPC Port
  - Database driver & path
  - JWT secret & expiry
  - Environment mode

#### 🐳 Docker Support
- `Dockerfile` dengan multi-stage build
- `docker-compose.yml` untuk easy deployment
- `.dockerignore` untuk optimasi build
- Health check endpoint `/health`
- Volume mounting untuk persistent data
- Environment variable configuration

### Changed
- Entity `User` sekarang memiliki field `Password`
- Repository menambahkan method `FindByEmail`
- Main.go diupdate untuk menggunakan config dan auth
- Database config sekarang menggunakan environment variables
- Semua user endpoints sekarang protected dengan JWT (kecuali auth endpoints)

### Documentation
- `AUTH.md` - Panduan lengkap authentication
- `DEPLOYMENT.md` - Panduan deployment dengan Docker
- `CHANGELOG.md` - Dokumentasi perubahan
- `Makefile` - Helper commands untuk development
- Updated `README.md` dengan fitur baru
- Postman collection baru dengan auth support

### Files Added
```
config/config.go
middleware/auth.go
middleware/grpc_auth.go
service/auth_service.go
controller/auth_controller.go
dto/auth_dto.go
.env.example
Dockerfile
docker-compose.yml
.dockerignore
Makefile
AUTH.md
DEPLOYMENT.md
CHANGELOG.md
User_CRUD_API_with_Auth.postman_collection.json
```

### Dependencies Added
- `github.com/golang-jwt/jwt/v5` - JWT token generation & validation
- `golang.org/x/crypto/bcrypt` - Password hashing

### Breaking Changes
⚠️ **IMPORTANT**: Ini adalah breaking change!

1. Semua endpoint `/users/*` sekarang memerlukan JWT token
2. Database schema berubah (menambahkan field `password`)
3. Existing users tidak akan bisa login (perlu re-register)

### Migration Guide

Jika sudah ada data user sebelumnya:

1. Backup database lama
2. Hapus database lama atau migrate manual
3. User perlu re-register dengan password

### Security Notes

- JWT_SECRET HARUS diset di production
- Default secret key hanya untuk development
- Gunakan HTTPS di production
- Password minimal 6 karakter (bisa ditingkatkan)

## [1.0.0] - Initial Release

- Basic CRUD operations
- REST API dengan Gin
- gRPC server
- SQLite database
- Clean architecture
- Unit tests untuk service & gRPC layer
