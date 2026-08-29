package controllers

import (
	"laundry-app/config"
	"laundry-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BayarPesanan - Membuat transaksi pembayaran baru
func BayarPesanan(c *gin.Context) {
	var input models.Pembayaran

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan ke tabel pembayaran
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pembayaran"})
		return
	}

	// Update status di tabel order menjadi 'Lunas' atau 'Selesai'
	config.DB.Model(&models.Order{}).Where("id_order = ?", input.IDOrder).Update("status", "Selesai")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pembayaran berhasil dilakukan!",
		"data":    input,
	})
}