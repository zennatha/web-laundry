package main

import (
	"laundry-app/config"
	"laundry-app/controllers"
	"laundry-app/middleware"
	"laundry-app/policy" // Import package policy

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
	r.POST("/api/login/admin", controllers.LoginAdmin) // <-- Route Login Admin Baru
	r.GET("/api/layanan", controllers.GetLayanan)

	// Route Terproteksi (Wajib membawa JWT Token)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Route Order (Bisa diakses Pelanggan & Admin)
		protected.POST("/orders", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.CreateOrder)
		protected.GET("/orders/pelanggan/:id_pelanggan", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.GetOrdersByPelanggan)
		protected.PUT("/orders/cancel/:id_order", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.CancelOrder)

		// Route khusus Update Status (HANYA Admin)
		protected.PUT("/orders/status/:id_order", middleware.RequireRole(policy.Admin), controllers.UpdateStatusOrder)

		// Route Pembayaran (Bisa diakses Pelanggan & Admin)
		protected.POST("/pembayaran", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.BayarPesanan)
	}

	r.Run(":8080")
}
