# Quick Start Guide

Panduan cepat untuk menjalankan aplikasi User CRUD API.

## Prerequisites

- Go 1.24+
- Docker & Docker Compose (opsional, untuk deployment)

## 1. Clone & Setup

```bash
git clone https://github.com/dwisusanto007/api-user-crud-go.git
cd api-user-crud-go
```

## 2. Install Dependencies

```bash
go mod download
```

## 3. Configuration (Opsional)

Untuk development, aplikasi akan menggunakan default config. Untuk production atau custom config:

```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

## 4. Run Aplikasi

### Cara 1: Langsung dengan Go

```bash
go run main.go
```

### Cara 2: Build dulu, lalu run

```bash
make build
./api-user-crud-go
```

### Cara 3: Dengan Docker

```bash
make docker-up
```

Aplikasi akan berjalan di:
- REST API: http://localhost:8080
- gRPC Server: localhost:50051

## 5. Test API

### A. Menggunakan cURL

#### Register user baru
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "age": 25
  }'
```

Response akan berisi token JWT:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25
  }
}
```

#### Simpan token untuk request selanjutnya
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### Get all users (dengan token)
```bash
curl http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN"
```

### B. Menggunakan Postman

1. Import file `User_CRUD_API_with_Auth.postman_collection.json`
2. Jalankan request "Register" atau "Login"
3. Token akan otomatis tersimpan di collection variable
4. Request lainnya akan otomatis menggunakan token tersebut

### C. Menggunakan grpcurl

```bash
# Install grpcurl
brew install grpcurl

# Register (public, tidak perlu token)
grpcurl -plaintext \
  -d '{"name":"John","email":"john@example.com","password":"password123","age":25}' \
  localhost:50051 user.UserService/Register

# Get all users (dengan token)
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN" \
  -d '{}' \
  localhost:50051 user.UserService/GetAllUsers
```

## 6. Run Tests

```bash
# Run all tests
make test

# Run tests dengan coverage
make test-coverage
```

## 7. Stop Aplikasi

### Jika run dengan Go
Tekan `Ctrl+C`

### Jika run dengan Docker
```bash
make docker-down
```

## Common Commands

```bash
make help           # Tampilkan semua available commands
make build          # Build aplikasi
make run            # Run aplikasi
make test           # Run tests
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-up      # Start dengan docker-compose
make docker-down    # Stop docker-compose
make docker-logs    # View docker logs
```

## Troubleshooting

### Port sudah digunakan
Jika port 8080 atau 50051 sudah digunakan, ubah di `.env`:
```
HTTP_PORT=8081
GRPC_PORT=50052
```

### Database error
Hapus file database dan restart:
```bash
rm test.db
go run main.go
```

### Docker build error
Pastikan Docker daemon berjalan:
```bash
docker ps
```

## Next Steps

- Baca [AUTH.md](AUTH.md) untuk detail authentication
- Baca [DEPLOYMENT.md](DEPLOYMENT.md) untuk deployment guide
- Baca [README.md](README.md) untuk dokumentasi lengkap
