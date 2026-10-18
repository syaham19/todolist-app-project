package main

import (
	"log"

	"github.com/gin-gonic/gin"

	// Import package dari project-mu
	// Ganti "todolist-project" sesuai nama module di go.mod kamu
	authHttp "todolist-project/internal/auth/delivery/http"
	authRepo "todolist-project/internal/auth/repository"
	authUseCase "todolist-project/internal/auth/usecase"
)

func main() {
	// 1. Inisialisasi Gin Router
	r := gin.Default()

	// 2. Inisialisasi Layer Repository (In-Memory)
	userRepository := authRepo.NewUserRepositoryMem()

	// 3. Inisialisasi Layer Usecase (Masukkan repository ke dalamnya)
	authUsecase := authUseCase.NewAuthUseCase(userRepository)

	// 4. Inisialisasi Layer Handler (Masukkan usecase ke dalamnya)
	authHandler := authHttp.NewAuthHandler(authUsecase)

	// 5. Daftarkan Route API Gin
	api := r.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.GET("/users", authHandler.GetAllUsers) // Tambahkan endpoint untuk mengambil semua user
	}

	// 6. Jalankan Server di Port 8080
	log.Println("Server running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
