package controller

import (
	"api-user-crud-go/dto"
	"api-user-crud-go/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController menangani HTTP requests untuk endpoint User.
// Controller menerima request, memanggil service, dan mengembalikan response.
type UserController struct {
	userService service.UserService
}

// NewUserController membuat instance baru UserController.
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser handler untuk POST /users - Membuat user baru.
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	// Bind dan validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}

	// Panggil service untuk membuat user
	user, err := ctrl.userService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create user",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUsers handler untuk GET /users - Mengambil semua user.
func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve users",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser handler untuk GET /users/:id - Mengambil user berdasarkan ID.
func (ctrl *UserController) GetUser(c *gin.Context) {
	// Parse ID dari parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
		})
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handler untuk PUT /users/:id - Mengupdate user.
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	// Parse ID dari parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}

	user, err := ctrl.userService.UpdateUser(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Failed to update user",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handler untuk DELETE /users/:id - Menghapus user.
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	// Parse ID dari parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
		})
		return
	}

	err = ctrl.userService.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Failed to delete user",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
