package main

import (
	"api-user-crud-go/config"
	"api-user-crud-go/controller"
	"api-user-crud-go/exception"
	"api-user-crud-go/repository"
	"api-user-crud-go/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// ==========================================
	// 1. INISIALISASI DATABASE
	// ==========================================
	db := config.InitDB()

	// ==========================================
	// 2. DEPENDENCY INJECTION (Wiring Layers)
	// ==========================================
	// Repository layer - mengakses database
	userRepo := repository.NewUserRepository(db)

	// Service layer - business logic, menggunakan repository
	userService := service.NewUserService(userRepo)

	// Controller layer - HTTP handlers, menggunakan service
	userController := controller.NewUserController(userService)

	// ==========================================
	// 3. SETUP GIN ROUTER & MIDDLEWARE
	// ==========================================
	router := gin.New()

	// Middleware global
	router.Use(exception.LoggerMiddleware()) // Logging setiap request
	router.Use(exception.Recovery())         // Recovery dari panic
	router.Use(exception.ErrorHandler())     // Handle error secara konsisten

	// ==========================================
	// 4. REGISTER ROUTES (API Endpoints)
	// ==========================================
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userController.CreateUser)       // POST /users
		userRoutes.GET("", userController.GetUsers)          // GET /users
		userRoutes.GET("/:id", userController.GetUser)       // GET /users/:id
		userRoutes.PUT("/:id", userController.UpdateUser)    // PUT /users/:id
		userRoutes.DELETE("/:id", userController.DeleteUser) // DELETE /users/:id
	}

	// ==========================================
	// 5. START SERVER
	// ==========================================
	log.Println("✓ Server berjalan di http://localhost:8080")
	log.Println("✓ API Endpoints:")
	log.Println("  - POST   /users")
	log.Println("  - GET    /users")
	log.Println("  - GET    /users/:id")
	log.Println("  - PUT    /users/:id")
	log.Println("  - DELETE /users/:id")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
