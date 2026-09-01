package controllers

import (
	"net/http"
	"os"
	"time"

	"laundry-app/config"
	"laundry-app/models"
	"laundry-app/policy"

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

	// Buat JWT Token dengan tambahan klaim role
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"id_pelanggan": pelanggan.IDPelanggan,
		"email":        pelanggan.Email,
		"role":         policy.Pelanggan, // Menambahkan role pelanggan
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id_pelanggan": pelanggan.IDPelanggan,
			"nama":         pelanggan.Nama,
			"email":        pelanggan.Email,
			"role":         policy.Pelanggan,
		},
	})
}

// Struct pendukung untuk Login Admin
type AdminModel struct {
	IDAdmin  uint   `gorm:"primaryKey;column:id_admin"`
	Nama     string `gorm:"column:nama"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (AdminModel) TableName() string {
	return "admin"
}

func LoginAdmin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin AdminModel
	if err := config.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password admin salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password admin salah"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"id_admin": admin.IDAdmin,
		"email":    admin.Email,
		"role":     policy.Admin, // Menambahkan role admin
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Admin berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id_admin": admin.IDAdmin,
			"nama":     admin.Nama,
			"email":    admin.Email,
			"role":     policy.Admin,
		},
	})
}
