# User CRUD API - Go

REST API sederhana untuk manajemen User menggunakan **Gin Gonic** dan **GORM** dengan arsitektur berlapis (layered architecture).

## 📋 Deskripsi

API ini adalah implementasi CRUD (Create, Read, Update, Delete) untuk entitas User dengan clean architecture yang memisahkan concern ke dalam beberapa layer:

- **Entity Layer**: Database models
- **DTO Layer**: Data Transfer Objects untuk request/response
- **Repository Layer**: Data access logic
- **Service Layer**: Business logic
- **Controller Layer**: HTTP handlers
- **Exception Layer**: Global error handling
- **Config Layer**: Database configuration

## 🚀 Tech Stack

- **Language**: Go 1.23+
- **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: SQLite (untuk development, mudah diganti ke PostgreSQL/MySQL)

## 📁 Project Structure

```
api-user-crud-go/
├── config/             # Database configuration
│   └── database.go
├── entity/             # Database models
│   └── user.go
├── dto/                # Data Transfer Objects
│   └── user_dto.go
├── repository/         # Data access layer
│   └── user_repository.go
├── service/            # Business logic layer
│   └── user_service.go
├── controller/         # HTTP handlers
│   └── user_controller.go
├── exception/          # Error handling middleware
│   └── error_handler.go
├── main.go            # Application entry point
├── go.mod
└── User_CRUD_API.postman_collection.json
```

## 🛠️ Installation

### Prerequisites

- Go 1.23 atau lebih tinggi
- Git

### Clone Repository

```bash
git clone <repository-url>
cd api-user-crud-go
```

### Install Dependencies

```bash
go mod tidy
```

## ▶️ Running the Application

```bash
# Jalankan aplikasi
go run main.go

# Atau build dan jalankan
go build .
./api-user-crud-go
```

Server akan berjalan di `http://localhost:8080`

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create new user |
| GET | `/users` | Get all users |
| GET | `/users/:id` | Get user by ID |
| PUT | `/users/:id` | Update user |
| DELETE | `/users/:id` | Delete user |

## 📝 API Usage Examples

### Create User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 25
}
```

### Get All Users

```bash
curl http://localhost:8080/users
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25
  }
]
```

### Get User by ID

```bash
curl http://localhost:8080/users/1
```

### Update User

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "age": 30
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/users/1
```

## 🧪 Testing with Postman

1. Import file `User_CRUD_API.postman_collection.json` ke Postman
2. Collection berisi semua endpoint yang siap digunakan
3. Jalankan request sesuai kebutuhan

## 🏗️ Architecture

API ini menggunakan **Layered Architecture** dengan dependency injection:

```
Database (GORM)
    ↓
Repository (data access)
    ↓
Service (business logic)
    ↓
Controller (HTTP handlers)
    ↓
Gin Router (routes + middleware)
```

### Keuntungan Arsitektur Ini:

- ✅ **Separation of Concerns**: Setiap layer punya tanggung jawab jelas
- ✅ **Testability**: Mudah di-unit test dengan mocking
- ✅ **Maintainability**: Mudah menemukan dan memodifikasi code
- ✅ **Scalability**: Mudah menambah fitur baru
- ✅ **Reusability**: Service & repository bisa dipakai ulang

## 📦 Database

Database menggunakan SQLite dengan file `test.db` yang akan dibuat otomatis saat aplikasi pertama kali dijalankan.

### Database Schema

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

Input validasi otomatis menggunakan Gin binding tags:

- **name**: Required
- **email**: Required, harus format email valid
- **age**: Required, minimal 1

Contoh error response:
```json
{
  "error": "Invalid input",
  "message": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

## 🚧 Future Enhancements

- [ ] Unit tests untuk setiap layer
- [ ] Pagination untuk GET /users
- [ ] Authentication & Authorization (JWT)
- [ ] Swagger/OpenAPI documentation
- [ ] Docker support
- [ ] CI/CD pipeline
- [ ] Logging dengan structured logger (zerolog/zap)
- [ ] Database migration files
- [ ] Environment-based configuration

## 📄 License

MIT License

## 👨‍💻 Author

Dibuat sebagai project pembelajaran Golang dengan Gin & GORM.
