package main

import (
	"api-user-crud-go/config"
	"api-user-crud-go/controller"
	"api-user-crud-go/exception"
	grpcserver "api-user-crud-go/grpcserver"
	"api-user-crud-go/proto"
	"api-user-crud-go/repository"
	"api-user-crud-go/service"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	// 3. START gRPC SERVER (port :50051)
	// ==========================================
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Gagal mengaktifkan gRPC listener: %v", err)
		}

		grpcServer := grpc.NewServer()

		// Register UserService gRPC handler (berbagi userService yang sama)
		proto.RegisterUserServiceServer(grpcServer, grpcserver.NewUserGRPCServer(userService))

		// Register reflection service (untuk grpcurl & tooling lainnya)
		reflection.Register(grpcServer)

		log.Println("✓ gRPC Server berjalan di grpc://localhost:50051")
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
	// 4. SETUP GIN ROUTER & MIDDLEWARE
	// ==========================================
	router := gin.New()

	// Middleware global
	router.Use(exception.LoggerMiddleware()) // Logging setiap request
	router.Use(exception.Recovery())         // Recovery dari panic
	router.Use(exception.ErrorHandler())     // Handle error secara konsisten

	// ==========================================
	// 5. REGISTER ROUTES (API Endpoints)
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
	// 6. START HTTP SERVER (port :8080)
	// ==========================================
	log.Println("✓ HTTP Server berjalan di http://localhost:8080")
	log.Println("✓ REST API Endpoints:")
	log.Println("  - POST   /users")
	log.Println("  - GET    /users")
	log.Println("  - GET    /users/:id")
	log.Println("  - PUT    /users/:id")
	log.Println("  - DELETE /users/:id")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan HTTP server:", err)
	}
}
