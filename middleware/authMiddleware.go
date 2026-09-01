package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, token tidak ditemukan!"})
			c.Abort()
			return
		}

		// Header format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah! Gunakan 'Bearer <token>'"})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("method signing tidak valid")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa!"})
			c.Abort()
			return
		}

		// Simpan klaim data user dan role dari token ke context Gin
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("id_pelanggan", claims["id_pelanggan"])
			c.Set("role", claims["role"]) // <-- Menambahkan role ke context Gin
		}

		c.Next()
	}
}

// Middleware baru untuk mengecek Policy / Role akses
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role tidak ditemukan pada token!"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Format role tidak valid!"})
			c.Abort()
			return
		}

		// Cek apakah role user ada di dalam daftar role yang diizinkan
		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak, Anda tidak memiliki izin untuk halaman ini!"})
		c.Abort()
	}
}
