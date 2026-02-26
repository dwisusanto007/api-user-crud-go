# User CRUD API - Go

REST & gRPC API untuk manajemen User menggunakan **Gin Gonic**, **GORM**, dan **gRPC** dengan arsitektur berlapis (layered architecture).

## 📋 Deskripsi

API ini adalah implementasi CRUD (Create, Read, Update, Delete) untuk entitas User dengan clean architecture yang memisahkan concern ke dalam beberapa layer:

- **Entity Layer**: Database models
- **DTO Layer**: Data Transfer Objects untuk request/response
- **Repository Layer**: Data access logic
- **Service Layer**: Business logic (shared antara REST & gRPC)
- **Controller Layer**: HTTP handlers (REST)
- **gRPC Server Layer**: gRPC handlers
- **Middleware Layer**: Authentication & authorization
- **Exception Layer**: Global error handling
- **Config Layer**: Environment-based configuration

## 🚀 Quick Start

Lihat [QUICKSTART.md](QUICKSTART.md) untuk panduan cepat menjalankan aplikasi.

## 🚀 Tech Stack

- **Language**: Go 1.24+
- **HTTP Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: SQLite
- **RPC**: [gRPC](https://grpc.io/) + [Protocol Buffers](https://protobuf.dev/)

## 📁 Project Structure

```
api-user-crud-go/
├── config/                 # Database configuration
│   └── database.go
├── entity/                 # Database models
│   └── user.go
├── dto/                    # Data Transfer Objects
│   └── user_dto.go
├── repository/             # Data access layer
│   └── user_repository.go
├── service/                # Business logic layer (shared REST & gRPC)
│   └── user_service.go
├── controller/             # REST HTTP handlers
│   └── user_controller.go
├── grpcserver/             # gRPC handlers
│   └── user_grpc_server.go
├── proto/                  # Protobuf definitions & generated code
│   ├── user.proto
│   ├── user.pb.go
│   └── user_grpc.pb.go
├── exception/              # Error handling middleware
│   └── error_handler.go
├── main.go                 # Application entry point
├── go.mod
└── User_CRUD_API.postman_collection.json
```

## 🛠️ Installation

### Prerequisites

- Go 1.24 atau lebih tinggi
- Git
- `protoc` (hanya jika ingin regenerate proto) → `brew install protobuf`

### Clone & Install

```bash
git clone https://github.com/dwisusanto007/api-user-crud-go.git
cd api-user-crud-go
go mod tidy
```

## ▶️ Running the Application

```bash
go run main.go
```

Server akan berjalan di dua port sekaligus:

| Server | Port | Protocol |
|--------|------|----------|
| REST API | `:8080` | HTTP/JSON |
| gRPC Server | `:50051` | HTTP/2 + Protobuf |

## 📡 REST API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create new user |
| GET | `/users` | Get all users |
| GET | `/users/:id` | Get user by ID |
| PUT | `/users/:id` | Update user |
| DELETE | `/users/:id` | Delete user |

### REST Usage Examples

```bash
# Create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "age": 25}'

# Get all users
curl http://localhost:8080/users

# Get user by ID
curl http://localhost:8080/users/1

# Update user
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "age": 30}'

# Delete user
curl -X DELETE http://localhost:8080/users/1
```

## ⚡ gRPC Endpoints

| RPC Method | Request | Response |
|---|---|---|
| `CreateUser` | `CreateUserRequest` | `UserMessage` |
| `GetAllUsers` | `GetAllUsersRequest` | `GetAllUsersResponse` |
| `GetUser` | `GetUserRequest` | `UserMessage` |
| `UpdateUser` | `UpdateUserRequest` | `UserMessage` |
| `DeleteUser` | `DeleteUserRequest` | `DeleteUserResponse` |

### gRPC Usage with grpcurl

Install grpcurl: `brew install grpcurl`

```bash
# Create user
grpcurl -plaintext -d '{"name":"John Doe","email":"john@example.com","age":25}' \
  localhost:50051 user.UserService/CreateUser

# Get all users
grpcurl -plaintext localhost:50051 user.UserService/GetAllUsers

# Get user by ID
grpcurl -plaintext -d '{"id":1}' localhost:50051 user.UserService/GetUser

# Update user
grpcurl -plaintext -d '{"id":1,"name":"Jane Doe","age":30}' \
  localhost:50051 user.UserService/UpdateUser

# Delete user
grpcurl -plaintext -d '{"id":1}' localhost:50051 user.UserService/DeleteUser
```

> Reflection service sudah diregistrasi — tidak perlu flag `--proto` saat menggunakan grpcurl.

## 🏗️ Architecture

```
HTTP Client (:8080)       gRPC Client (:50051)
        |                         |
   Gin Controller           gRPC Server
        |                         |
        +--------> UserService <--+
                       |
                  Repository
                       |
                  SQLite DB
```

### Keuntungan Arsitektur Ini:

- ✅ **Separation of Concerns**: Setiap layer punya tanggung jawab jelas
- ✅ **Shared Business Logic**: REST & gRPC menggunakan `UserService` yang sama
- ✅ **Testability**: Mudah di-unit test dengan mocking
- ✅ **Scalability**: Mudah menambah fitur baru di kedua protokol

## 📦 Database

SQLite dengan file `test.db` yang dibuat otomatis saat aplikasi pertama dijalankan.

**Table: users**

| Column | Type | Constraint |
|--------|------|------------|
| id | INTEGER | PRIMARY KEY |
| name | TEXT | NOT NULL |
| email | TEXT | UNIQUE, NOT NULL |
| age | INTEGER | |
| created_at | DATETIME | |
| updated_at | DATETIME | |
| deleted_at | DATETIME | (soft delete) |

## 🔒 Validasi

Input validasi otomatis:

- **name**: Required
- **email**: Required, format email valid
- **age**: Required, minimal 1

## 🧪 Testing with Postman

1. Import `User_CRUD_API.postman_collection.json` ke Postman
2. Collection berisi semua REST endpoint yang siap digunakan

## 🚧 Future Enhancements

- [x] Unit tests untuk setiap layer
- [x] Authentication & Authorization (JWT / gRPC interceptor)
- [x] Docker support
- [x] Environment-based configuration
- [ ] Pagination untuk GetAllUsers
- [ ] Swagger/OpenAPI documentation
- [ ] CI/CD pipeline
- [ ] gRPC streaming endpoints

## 📚 Additional Documentation

- [Authentication Guide](AUTH.md) - Panduan lengkap JWT authentication
- [Deployment Guide](DEPLOYMENT.md) - Panduan deployment dengan Docker
- [Security Guidelines](SECURITY.md) - Security best practices dan checklist
- [Quick Start](QUICKSTART.md) - Quick start guide

## 🔐 Authentication

API ini menggunakan JWT untuk autentikasi. Lihat [AUTH.md](AUTH.md) untuk detail lengkap.

### Quick Start

1. Register user baru:
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"password123","age":25}'
```

2. Login dan dapatkan token:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```

3. Gunakan token untuk akses protected endpoints:
```bash
curl http://localhost:8080/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🐳 Docker Deployment

```bash
# Build dan run dengan docker-compose
docker-compose up -d

# Atau build manual
docker build -t user-crud-api .
docker run -p 8080:8080 -p 50051:50051 user-crud-api
```

Lihat [DEPLOYMENT.md](DEPLOYMENT.md) untuk panduan lengkap.

## ⚙️ Configuration

Copy `.env.example` ke `.env` dan sesuaikan konfigurasi:

```bash
cp .env.example .env
```

Environment variables yang tersedia:
- `HTTP_PORT` - Port REST API (default: 8080)
- `GRPC_PORT` - Port gRPC (default: 50051)
- `JWT_SECRET` - Secret key untuk JWT (WAJIB di production)
- `JWT_EXPIRY_HOURS` - Durasi token (default: 24 jam)
- `ENV` - Environment: development/production

## 📄 License

MIT License

## 👨‍💻 Author

Dibuat sebagai project pembelajaran Golang dengan Gin, GORM, dan gRPC.
