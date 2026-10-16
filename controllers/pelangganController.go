package controllers

import (
	"net/http"

	"laundry-app/config"
	"laundry-app/models"

	"github.com/gin-gonic/gin"
)

// GET /api/pelanggan
func GetPelanggan(c *gin.Context) {
	var pelanggan []models.Pelanggan

	result := config.DB.Find(&pelanggan)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data pelanggan",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pelanggan berhasil diambil",
		"data":    pelanggan,
	})
}

// GET /api/pelanggan/:id
func GetPelangganByID(c *gin.Context) {
	var pelanggan models.Pelanggan

	id := c.Param("id")

	result := config.DB.First(&pelanggan, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pelanggan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pelanggan ditemukan",
		"data":    pelanggan,
	})
}

// POST /api/pelanggan
func CreatePelanggan(c *gin.Context) {
	var pelanggan models.Pelanggan

	if err := c.ShouldBindJSON(&pelanggan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Data pelanggan tidak valid",
			"error":   err.Error(),
		})
		return
	}

	result := config.DB.Create(&pelanggan)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menambahkan pelanggan",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pelanggan berhasil ditambahkan",
		"data":    pelanggan,
	})
}

// PUT /api/pelanggan/:id
func UpdatePelanggan(c *gin.Context) {
	id := c.Param("id")

	var pelanggan models.Pelanggan

	result := config.DB.First(&pelanggan, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pelanggan tidak ditemukan",
		})
		return
	}

	var input struct {
		Nama   string `json:"nama"`
		NoHp   string `json:"no_hp"`
		Alamat string `json:"alamat"`
		Email  string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	pelanggan.Nama = input.Nama
	pelanggan.NoHp = input.NoHp
	pelanggan.Alamat = input.Alamat
	pelanggan.Email = input.Email

	result = config.DB.Save(&pelanggan)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengubah data pelanggan",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pelanggan berhasil diubah",
		"data":    pelanggan,
	})
}

// DELETE /api/pelanggan/:id
func DeletePelanggan(c *gin.Context) {
	id := c.Param("id")

	var pelanggan models.Pelanggan

	result := config.DB.First(&pelanggan, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pelanggan tidak ditemukan",
		})
		return
	}

	result = config.DB.Delete(&pelanggan)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus pelanggan",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pelanggan berhasil dihapus",
	})
}
