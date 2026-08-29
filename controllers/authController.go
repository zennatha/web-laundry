package controllers

import (
	"net/http"
	"os"
	"time"

	"laundry-app/config"
	"laundry-app/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Nama     string `json:"nama" binding:"required"`
		NoHp     string `json:"no_hp" binding:"required"`
		Alamat   string `json:"alamat"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Cek dulu apakah email sudah terdaftar di database
	var existingUser models.Pelanggan
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar!"})
		return
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	pelanggan := models.Pelanggan{
		Nama:     input.Nama,
		NoHp:     input.NoHp,
		Alamat:   input.Alamat,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// 3. Simpan ke Database
	if err := config.DB.Create(&pelanggan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id_pelanggan": pelanggan.IDPelanggan,
			"nama":         pelanggan.Nama,
			"email":        pelanggan.Email,
		},
	})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pelanggan models.Pelanggan
	if err := config.DB.Where("email = ?", input.Email).First(&pelanggan).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pelanggan.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// --- PERUBAHAN JWT DI SINI ---
	// Buat JWT Token (Berlaku 24 Jam)
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"id_pelanggan": pelanggan.IDPelanggan,
		"email":        pelanggan.Email,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	// Kembalikan token ke response JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id_pelanggan": pelanggan.IDPelanggan,
			"nama":         pelanggan.Nama,
			"email":        pelanggan.Email,
		},
	})
}
