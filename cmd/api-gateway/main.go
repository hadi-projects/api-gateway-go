// cmd/api-gateway/main.go
package main

import (
	"log"
	"os"

	"api-gateway-go/pkg/config"
	"api-gateway-go/pkg/database"
	"api-gateway-go/pkg/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Impor GORM
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// Inisialisasi Database GORM
	var db *gorm.DB // Sekarang bertipe *gorm.DB
	db, err = database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database GORM: %v", err)
	}
	// GORM tidak memiliki fungsi Close() secara langsung pada *gorm.DB,
	// tapi Anda bisa mendapatkan *sql.DB underlying jika perlu menutupnya secara eksplisit.
	// Namun, biasanya koneksi dikelola oleh GORM dan akan tertutup saat aplikasi berhenti.
	// Untuk penutupan eksplisit (misalnya, dalam defer setelah aplikasi selesai):
	sqlDB, err := db.DB()
	if err == nil { // Hanya jika berhasil mendapatkan *sql.DB
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Gagal menutup koneksi database (underlying sql.DB): %v", err)
			} else {
				log.Println("Koneksi database (underlying sql.DB) berhasil ditutup.")
			}
		}()
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	// Setup Rute, sekarang teruskan *gorm.DB
	routes.SetupRoutes(router, cfg, db)

	port := cfg.ServerPort
	if port == "" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			port = envPort
		} else {
			port = "8080"
		}
	}

	log.Printf("API Gateway (dengan GORM) siap dijalankan di port :%s", port)
	// ... log lainnya
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server Gin: %v", err)
	}
}
