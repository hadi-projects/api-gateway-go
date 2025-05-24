// pkg/handlers/auth_handler.go
package handlers

import (
	"net/http"
	"time"

	"api-gateway-go/pkg/config" // Sesuaikan dengan nama modul Anda

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Credentials adalah struktur untuk data login yang diterima
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Claims adalah struktur untuk data yang akan disimpan dalam token JWT
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// LoginHandler menangani permintaan login dan menghasilkan token JWT
func LoginHandler(appConfig config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds Credentials
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
			return
		}

		// --- MULAI VALIDASI KREDENSIAL (SIMULASI) ---
		// Di aplikasi nyata, Anda akan memvalidasi username dan password dengan database.
		// Untuk contoh ini, kita akan hardcode kredensial yang valid.
		var userID string
		if creds.Username == "user123" && creds.Password == "password123" {
			userID = "USR_001" // ID pengguna dari database
		} else if creds.Username == "admin" && creds.Password == "adminpass" {
			userID = "ADM_001"
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		// --- SELESAI VALIDASI KREDENSIAL (SIMULASI) ---

		// Tentukan waktu kedaluwarsa token (misalnya, 1 jam)
		expirationTime := time.Now().Add(1 * time.Hour)

		// Buat claims
		claims := &Claims{
			UserID:   userID,
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "api-gateway", // Nama penerbit token
				Subject:   userID,        // Subjek token (seringkali ID pengguna)
			},
		}

		// Buat token baru dengan metode signing HS256 dan claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Dapatkan string token yang ditandatangani menggunakan secret dari config
		tokenString, err := token.SignedString([]byte(appConfig.AuthSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":      tokenString,
			"expires_at": expirationTime.Format(time.RFC3339),
			"user_id":    userID,
			"username":   creds.Username,
		})
	}
}
