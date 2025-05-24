// cmd/api-gateway/main.go
package main

import (
	"log"
	"os"

	"api-gateway-go/pkg/config" // Sesuaikan dengan nama modul Anda
	"api-gateway-go/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Muat Konfigurasi
	// Path "." berarti mencari config.yaml di direktori yang sama dengan executable
	// atau di direktori kerja saat menjalankan `go run`
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// Inisialisasi Gin Engine
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("Menjalankan dalam mode PRODUKSI")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Println("Menjalankan dalam mode DEBUG")
	}
	// gin.Default() sudah menyertakan middleware Logger dan Recovery
	router := gin.Default()

	// Setup Rute dari pkg/routes
	routes.SetupRoutes(router, cfg)

	// Dapatkan port dari konfigurasi atau environment variable (untuk platform seperti Heroku)
	port := cfg.ServerPort
	if port == "" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			port = envPort
		} else {
			port = "8080" // Port default jika tidak ada yang diset
		}
	}

	log.Printf("API Gateway siap dijalankan di port :%s", port)
	log.Printf("Endpoint Health Check: http://localhost:%s/api/public/health", port)
	log.Printf("Secret Otentikasi (dummy): %s", cfg.AuthSecret) // Jangan log secret asli di produksi
	log.Printf("User Service Endpoint: %s", cfg.ServiceEndpoints["user_service"])
	log.Printf("Product Service Endpoint: %s", cfg.ServiceEndpoints["product_service"])

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server Gin: %v", err)
	}
}
