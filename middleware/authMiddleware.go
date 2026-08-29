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

		// Simpan klaim data user dari token ke context Gin
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("id_pelanggan", claims["id_pelanggan"])
		}

		c.Next()
	}
}
