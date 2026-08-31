package controllers

import (
	"net/http"

	"laundry-app/config"
	"laundry-app/models"

	"github.com/gin-gonic/gin"
)

// Membuat Pesanan Baru (Membuat Pesanan - Pelanggan/Admin)
func CreateOrder(c *gin.Context) {
	var input struct {
		IDPelanggan uint    `json:"id_pelanggan" binding:"required"`
		Layanan     string  `json:"layanan" binding:"required"`
		Berat       float64 `json:"berat" binding:"required"`
		Harga       int     `json:"harga" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		IDPelanggan: input.IDPelanggan,
		Layanan:     input.Layanan,
		Berat:       input.Berat,
		Harga:       input.Harga,
		Status:      "Proses",
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

// Melihat Status & Riwayat Laundry Pelanggan
func GetOrdersByPelanggan(c *gin.Context) {
	idPelanggan := c.Param("id_pelanggan")
	var orders []models.Order

	if err := config.DB.Where("id_pelanggan = ?", idPelanggan).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// Ubah Status Laundry (Khusus Admin)
func UpdateStatusOrder(c *gin.Context) {
	idOrder := c.Param("id_order")
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, idOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	config.DB.Model(&order).Update("status", input.Status)

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pesanan berhasil diperbarui",
		"data":    order,
	})
}

// Membatalkan Pesanan (Khusus Pelanggan/Admin)
func CancelOrder(c *gin.Context) {
	idOrder := c.Param("id_order")
	var order models.Order

	// 1. Cari pesanan berdasarkan ID
	if err := config.DB.First(&order, idOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	// 2. Cek apakah pesanan masih berstatus 'Proses' (hanya bisa dibatalkan jika belum selesai)
	if order.Status != "Proses" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan yang sudah selesai atau dibatalkan tidak dapat diubah"})
		return
	}

	// 3. Update status pesanan menjadi 'Canceled'
	config.DB.Model(&order).Update("status", "Canceled")

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dibatalkan",
		"data":    order,
	})
}
