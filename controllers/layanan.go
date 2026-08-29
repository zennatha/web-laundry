package controllers

import (
	"laundry-app/config"
	"laundry-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLayanan(c *gin.Context) {
	var daftarLayanan []models.Layanan

	if err := config.DB.Find(&daftarLayanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data layanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data layanan",
		"data":    daftarLayanan,
	})
}
