package controllers

import (
	"net/http"

	"laundry-app/config"
	"laundry-app/models"

	"github.com/gin-gonic/gin"
)

// Membuat Pesanan Baru (Pelanggan / Admin)
func CreateOrder(c *gin.Context) {
	var input struct {
		IDPelanggan uint    `json:"id_pelanggan" binding:"required"`
		IDLayanan   uint    `json:"id_layanan"`
		Layanan     string  `json:"layanan"`
		TotalBerat  float64 `json:"total_berat"`
		Berat       float64 `json:"berat"`
		TotalHarga  int     `json:"total_harga"`
		Harga       int     `json:"harga"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalisasi field agar menerima request dari Vue frontend maupun Postman
	berat := input.TotalBerat
	if berat == 0 {
		berat = input.Berat
	}

	harga := input.TotalHarga
	if harga == 0 {
		harga = input.Harga
	}

	layanan := input.Layanan
	if layanan == "" && input.IDLayanan != 0 {
		var l models.Layanan
		if err := config.DB.First(&l, input.IDLayanan).Error; err == nil {
			layanan = l.NamaLayanan
		} else {
			layanan = "Cuci Reguler"
		}
	} else if layanan == "" {
		layanan = "Cuci Reguler"
	}

	order := models.Order{
		IDPelanggan: input.IDPelanggan,
		Layanan:     layanan,
		Berat:       berat,
		Harga:       harga,
		Status:      "Menunggu",
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

// Melihat Semua Pesanan (Khusus Admin)
func GetAllOrders(c *gin.Context) {
	var orders []models.Order

	if err := config.DB.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil seluruh data pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// Melihat Status & Riwayat Laundry berdasarkan ID Pelanggan
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
		StatusOrder string `json:"status_order"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusBaru := input.StatusOrder
	if statusBaru == "" {
		statusBaru = input.Status
	}

	if statusBaru == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak boleh kosong"})
		return
	}

	var order models.Order
	if err := config.DB.Where("id_order = ?", idOrder).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	if err := config.DB.Model(&order).Where("id_order = ?", idOrder).Update("status", statusBaru).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pesanan berhasil diperbarui",
		"data":    order,
	})
}

// Membatalkan Pesanan (Khusus Pelanggan / Admin)
func CancelOrder(c *gin.Context) {
	idOrder := c.Param("id_order")
	var order models.Order

	if err := config.DB.Where("id_order = ?", idOrder).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	if order.Status == "Selesai" || order.Status == "Diambil" || order.Status == "Canceled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan yang sudah selesai/dibatalkan tidak dapat diubah"})
		return
	}

	config.DB.Model(&order).Where("id_order = ?", idOrder).Update("status", "Canceled")

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dibatalkan",
		"data":    order,
	})
}
