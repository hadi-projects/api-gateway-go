// pkg/config/config.go
package config

import (
	"github.com/spf13/viper"
)

// Config menyimpan semua konfigurasi aplikasi.
type Config struct {
	ServerPort       string            `mapstructure:"SERVER_PORT"`
	AppEnv           string            `mapstructure:"APP_ENV"`
	AuthSecret       string            `mapstructure:"AUTH_SECRET"`
	ServiceEndpoints map[string]string `mapstructure:"SERVICE_ENDPOINTS"`
	RateLimit        RateLimitConfig   `mapstructure:"RATE_LIMIT"`
}

type RateLimitConfig struct {
	Enabled   bool `mapstructure:"ENABLED"`
	Requests  int  `mapstructure:"REQUESTS"`
	WindowSec int  `mapstructure:"WINDOW_SEC"`
}

// LoadConfig membaca konfigurasi dari file atau variabel environment.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)     // Path ke direktori tempat file config berada (root proyek)
	viper.SetConfigName("config") // Nama file config (tanpa ekstensi)
	viper.SetConfigType("yaml")   // Tipe file config

	viper.AutomaticEnv() // Baca variabel environment yang cocok

	// Set default values
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("AUTH_SECRET", "your-default-secret-key") // Ganti ini di produksi
	viper.SetDefault("RATE_LIMIT.ENABLED", true)
	viper.SetDefault("RATE_LIMIT.REQUESTS", 100)  // 100 requests
	viper.SetDefault("RATE_LIMIT.WINDOW_SEC", 60) // per 60 detik (1 menit)
	viper.SetDefault("SERVICE_ENDPOINTS.user_service", "http://localhost:8081")
	viper.SetDefault("SERVICE_ENDPOINTS.product_service", "http://localhost:8082")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file ditemukan tapi ada error lain saat parsing
			return
		}
		// Config file tidak ditemukan, tidak apa-apa, akan menggunakan default atau env var
	}

	err = viper.Unmarshal(&config)
	return
}
