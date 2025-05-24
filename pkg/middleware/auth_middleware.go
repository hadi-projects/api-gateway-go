// pkg/middleware/auth_middleware.go
package middleware

import (
	"api-gateway-go/pkg/handlers" // Untuk akses ke struct Claims
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware membuat middleware untuk otentikasi menggunakan JWT.
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims := &handlers.Claims{} // Menggunakan struct Claims dari handlers

		// Parse token JWT
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing adalah yang diharapkan (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			// Tangani error kedaluwarsa secara spesifik
			// if verr, ok := err.(*jwt.ValidationError); ok {
			// 	if verr.Errors&jwt.ValidationErrorMalformed != 0 {
			// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Malformed token"})
			// 		return
			// 	} else if verr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// 		// Token kedaluwarsa atau belum valid
			// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired or not valid yet"})
			// 		return
			// 	}
			// }
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Token valid. Anda bisa menyimpan informasi dari claims ke context jika perlu.
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		// c.Set("claims", claims) // Atau simpan seluruh claims

		log.Printf("Authenticated UserID: %s, Username: %s, ExpiresAt: %s",
			claims.UserID,
			claims.Username,
			claims.ExpiresAt.Format(time.RFC3339),
		)

		c.Next()
	}
}
