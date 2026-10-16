package controllers

import (
	"net/http"

	"laundry-app/config"
	"laundry-app/models"

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

	// Update status di tabel order menjadi 'Selesai'
	config.DB.Model(&models.Order{}).Where("id_order = ?", input.IDOrder).Update("status", "Selesai")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pembayaran berhasil dilakukan!",
		"data":    input,
	})
}

// GetPembayaran - Mengambil seluruh data pembayaran untuk Admin & Pelanggan
func GetPembayaran(c *gin.Context) {
	var listPembayaran []models.Pembayaran

	if err := config.DB.Find(&listPembayaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": listPembayaran})
}
