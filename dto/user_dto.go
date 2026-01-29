package dto

// CreateUserRequest adalah DTO untuk membuat user baru.
// Digunakan untuk menerima input dari POST /users.
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=1"`
}

// UpdateUserRequest adalah DTO untuk mengupdate user.
// Digunakan untuk menerima input dari PUT /users/:id.
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   int    `json:"age" binding:"omitempty,min=1"`
}

// UserResponse adalah DTO untuk response user.
// Digunakan untuk mengembalikan data user ke client.
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// ErrorResponse adalah DTO untuk response error.
// Digunakan untuk error handling yang konsisten.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
