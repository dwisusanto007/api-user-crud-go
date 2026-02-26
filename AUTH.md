# Authentication & Authorization Guide

API ini menggunakan JWT (JSON Web Token) untuk autentikasi.

## Flow Autentikasi

1. User melakukan registrasi atau login
2. Server mengembalikan JWT token
3. Client menyimpan token (localStorage, cookie, dll)
4. Client mengirim token di header `Authorization` untuk setiap request ke protected endpoints

## Public Endpoints (Tidak Perlu Token)

- `POST /auth/register` - Registrasi user baru
- `POST /auth/login` - Login user
- `GET /health` - Health check

## Protected Endpoints (Perlu Token)

Semua endpoint `/users/*` memerlukan JWT token:
- `POST /users` - Create user
- `GET /users` - Get all users
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

## REST API Examples

### 1. Register User Baru

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

Response:
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

### 2. Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response sama seperti register.

### 3. Akses Protected Endpoint

```bash
# Simpan token ke variable
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get all users (dengan token)
curl http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN"

# Create user (dengan token)
curl -X POST http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "age": 30
  }'
```

## gRPC Authentication

Untuk gRPC, token dikirim melalui metadata:

```bash
# Set token di metadata
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN" \
  -d '{}' \
  localhost:50051 user.UserService/GetAllUsers
```

## Token Information

- Token berlaku selama 24 jam (default, bisa diubah via `JWT_EXPIRY_HOURS`)
- Token berisi: `user_id`, `email`, `issued_at`, `expires_at`
- Token di-sign dengan `JWT_SECRET` (harus dijaga kerahasiaannya)

## Error Responses

### 401 Unauthorized

Token tidak valid, expired, atau tidak diberikan:

```json
{
  "error": "Invalid or expired token"
}
```

### 400 Bad Request

Validasi gagal (email sudah terdaftar, password kurang dari 6 karakter, dll):

```json
{
  "error": "email already registered"
}
```

## Security Best Practices

1. **Jangan hardcode JWT_SECRET** - Gunakan environment variable
2. **HTTPS di Production** - Selalu gunakan HTTPS untuk mencegah token dicuri
3. **Token Storage** - Simpan token dengan aman di client (HttpOnly cookies lebih aman dari localStorage)
4. **Token Expiry** - Set expiry time yang reasonable (24 jam default)
5. **Refresh Token** - Implementasi refresh token untuk UX yang lebih baik (future enhancement)
6. **Password Policy** - Minimal 6 karakter (bisa ditingkatkan)
7. **Rate Limiting** - Implementasi rate limiting untuk mencegah brute force attack
