package main

import (
	"laundry-app/config"
	"laundry-app/controllers"
	"laundry-app/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Hubungkan ke database MariaDB
	config.ConnectDatabase()

	// 2. Jalankan seeder admin (Otomatis hash password admin di DB jika belum ada)
	config.SeedAdmin()

	r := gin.Default()

	// Route Publik (Bisa diakses tanpa token)
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)
	r.GET("/api/layanan", controllers.GetLayanan)

	// Route Terproteksi (Wajib membawa JWT Token)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Route Order
		protected.POST("/orders", controllers.CreateOrder)
		protected.GET("/orders/pelanggan/:id_pelanggan", controllers.GetOrdersByPelanggan)
		protected.PUT("/orders/status/:id_order", controllers.UpdateStatusOrder)
		protected.PUT("/orders/cancel/:id_order", controllers.CancelOrder)

		// Route Pembayaran
		protected.POST("/pembayaran", controllers.BayarPesanan)
	}

	r.Run(":8080")
}
