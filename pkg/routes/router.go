// pkg/routes/router.go
package routes

import (
	"api-gateway-go/pkg/config" // Sesuaikan dengan nama modul Anda
	"api-gateway-go/pkg/handlers"
	"api-gateway-go/pkg/middleware"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRoutes(router *gin.Engine, cfg config.Config) {
	// Middleware Global
	router.Use(middleware.LoggingMiddleware())

	// CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // HATI-HATI: Untuk produksi, batasi origin
	// corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://yourfrontend.com"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Rate Limiting Per IP (Contoh: 5 request per detik, dengan burst 10)
	if cfg.RateLimit.Enabled {
		// Konversi dari config ke rate.Limit
		// rate.Limit adalah float64 yang merepresentasikan event per detik.
		// Jika cfg.RateLimit.Requests = 100 dan cfg.RateLimit.WindowSec = 60,
		// maka limitnya adalah 100/60 event/detik.
		// Burst adalah cfg.RateLimit.Requests
		limit := rate.Limit(float64(cfg.RateLimit.Requests) / float64(cfg.RateLimit.WindowSec))
		burst := cfg.RateLimit.Requests
		router.Use(middleware.RateLimitMiddlewarePerIP(limit, burst))
		log.Printf("Rate limiting enabled: %.2f req/sec, burst %d per IP", limit, burst)
	}

	// Public Routes
	public := router.Group("/api/public")
	{
		public.GET("/health", handlers.HealthCheck)
	}

	// Authentication Route
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handlers.LoginHandler(cfg)) // Mengirim config ke handler jika diperlukan
	}

	// API v1 Group
	apiV1 := router.Group("/api/v1")
	{
		// Service Pengguna
		userServiceTarget, ok := cfg.ServiceEndpoints["user_service"]
		if !ok || userServiceTarget == "" {
			log.Fatalf("URL untuk user_service tidak ditemukan di konfigurasi")
		}
		userServiceURL, err := url.Parse(userServiceTarget)
		if err != nil {
			log.Fatalf("URL user_service tidak valid: %v", err)
		}
		userProxy := handlers.NewProxyHandler(userServiceURL)

		userRoutes := apiV1.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware(cfg.AuthSecret)) // Semua endpoint user butuh auth
		{
			// Path /*proxyPath akan menangkap semua sub-path
			// Contoh: /api/v1/users/profile -> proxyPath = /profile
			// Contoh: /api/v1/users/123/orders -> proxyPath = /123/orders
			userRoutes.Any("/*proxyPath", userProxy.Handle)
		}

		// Service Produk
		productServiceTarget, ok := cfg.ServiceEndpoints["product_service"]
		if !ok || productServiceTarget == "" {
			log.Fatalf("URL untuk product_service tidak ditemukan di konfigurasi")
		}
		productServiceURL, err := url.Parse(productServiceTarget)
		if err != nil {
			log.Fatalf("URL product_service tidak valid: %v", err)
		}
		productProxy := handlers.NewProxyHandler(productServiceURL)

		productRoutes := apiV1.Group("/products")
		{
			// GET produk bisa publik
			productRoutes.GET("/*proxyPath", productProxy.Handle)

			// POST, PUT, DELETE produk butuh auth
			productProtected := productRoutes.Group("") // Grup kosong untuk menerapkan middleware tambahan
			productProtected.Use(middleware.AuthMiddleware(cfg.AuthSecret))
			{
				productProtected.POST("/*proxyPath", productProxy.Handle)
				productProtected.PUT("/*proxyPath", productProxy.Handle)
				productProtected.DELETE("/*proxyPath", productProxy.Handle)
			}
		}

		// Tambahkan service lain di sini jika ada
		// Contoh: Order Service (semua butuh auth)
		orderServiceTarget, ok := cfg.ServiceEndpoints["order_service"]
		if ok && orderServiceTarget != "" { // Hanya jika dikonfigurasi
			orderServiceURL, err := url.Parse(orderServiceTarget)
			if err != nil {
				log.Printf("Peringatan: URL order_service tidak valid: %v", err)
			} else {
				orderProxy := handlers.NewProxyHandler(orderServiceURL)
				orderRoutes := apiV1.Group("/orders")
				orderRoutes.Use(middleware.AuthMiddleware(cfg.AuthSecret))
				{
					orderRoutes.Any("/*proxyPath", orderProxy.Handle)
				}
			}
		}
	}

	// Fallback untuk rute yang tidak ditemukan
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "ROUTE_NOT_FOUND", "message": "Endpoint tidak ditemukan."})
	})
}
