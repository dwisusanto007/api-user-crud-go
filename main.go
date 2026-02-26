package main

import (
	"api-user-crud-go/config"
	"api-user-crud-go/controller"
	"api-user-crud-go/exception"
	grpcserver "api-user-crud-go/grpcserver"
	"api-user-crud-go/middleware"
	"api-user-crud-go/proto"
	"api-user-crud-go/repository"
	"api-user-crud-go/service"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// ==========================================
	// 1. LOAD CONFIGURATION
	// ==========================================
	cfg := config.LoadConfig()
	cfg.ValidateConfig()

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// ==========================================
	// 2. INISIALISASI DATABASE
	// ==========================================
	db := config.InitDB()

	// ==========================================
	// 3. DEPENDENCY INJECTION (Wiring Layers)
	// ==========================================
	// Repository layer - mengakses database
	userRepo := repository.NewUserRepository(db)

	// Service layer - business logic, menggunakan repository
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg)

	// Controller layer - HTTP handlers, menggunakan service
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)

	// ==========================================
	// 4. START gRPC SERVER (with auth interceptor)
	// ==========================================
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("Gagal mengaktifkan gRPC listener: %v", err)
		}

		// Create gRPC server with auth interceptor
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(middleware.GRPCAuthInterceptor(cfg)),
		)

		// Register UserService gRPC handler (berbagi userService yang sama)
		proto.RegisterUserServiceServer(grpcServer, grpcserver.NewUserGRPCServer(userService))

		// Register reflection service (untuk grpcurl & tooling lainnya)
		reflection.Register(grpcServer)

		log.Printf("✓ gRPC Server berjalan di grpc://localhost:%s\n", cfg.GRPCPort)
		log.Println("✓ gRPC Methods:")
		log.Println("  - UserService/CreateUser")
		log.Println("  - UserService/GetAllUsers")
		log.Println("  - UserService/GetUser")
		log.Println("  - UserService/UpdateUser")
		log.Println("  - UserService/DeleteUser")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Gagal menjalankan gRPC server: %v", err)
		}
	}()

	// ==========================================
	// 5. SETUP GIN ROUTER & MIDDLEWARE
	// ==========================================
	router := gin.New()

	// Middleware global
	router.Use(exception.LoggerMiddleware()) // Logging setiap request
	router.Use(exception.Recovery())         // Recovery dari panic
	router.Use(exception.ErrorHandler())     // Handle error secara konsisten

	// ==========================================
	// 6. REGISTER ROUTES (API Endpoints)
	// ==========================================
	
	// Health check endpoint (public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register) // POST /auth/register
		authRoutes.POST("/login", authController.Login)       // POST /auth/login
	}

	// User routes (protected with JWT)
	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.JWTAuth(cfg)) // Apply JWT middleware
	{
		userRoutes.POST("", userController.CreateUser)       // POST /users
		userRoutes.GET("", userController.GetUsers)          // GET /users
		userRoutes.GET("/:id", userController.GetUser)       // GET /users/:id
		userRoutes.PUT("/:id", userController.UpdateUser)    // PUT /users/:id
		userRoutes.DELETE("/:id", userController.DeleteUser) // DELETE /users/:id
	}

	// ==========================================
	// 7. START HTTP SERVER
	// ==========================================
	log.Printf("✓ HTTP Server berjalan di http://localhost:%s\n", cfg.HTTPPort)
	log.Println("✓ REST API Endpoints:")
	log.Println("  Public:")
	log.Println("    - GET    /health")
	log.Println("    - POST   /auth/register")
	log.Println("    - POST   /auth/login")
	log.Println("  Protected (requires JWT):")
	log.Println("    - POST   /users")
	log.Println("    - GET    /users")
	log.Println("    - GET    /users/:id")
	log.Println("    - PUT    /users/:id")
	log.Println("    - DELETE /users/:id")

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("Gagal menjalankan HTTP server:", err)
	}
}
