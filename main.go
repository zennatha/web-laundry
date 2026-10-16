package main

import (
	"time"

	"laundry-app/config"
	"laundry-app/controllers"
	"laundry-app/middleware"
	"laundry-app/policy"

	"github.com/gin-contrib/cors" // Import middleware CORS
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Hubungkan ke database MariaDB
	config.ConnectDatabase()

	// 2. Jalankan seeder admin
	config.SeedAdmin()

	r := gin.Default()

	// 3. Konfigurasi CORS Lengkap & Valid
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Mengizinkan semua origin/IP Laptop B secara dinamis tanpa bentrok AllowCredentials
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Menangani preflight request OPTIONS secara langsung
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})

	// Route Publik (Bisa diakses tanpa token)
	r.POST("/api/register", controllers.Register)
	r.POST("/api/register/admin", controllers.RegisterAdmin)
	r.POST("/api/login", controllers.Login)
	r.POST("/api/login/admin", controllers.LoginAdmin)
	r.GET("/api/layanan", controllers.GetLayanan)

	// Route Terproteksi (Wajib membawa JWT Token)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Route Pelanggan (Admin & Pelanggan)
		protected.GET("/pelanggan", middleware.RequireRole(policy.Admin), controllers.GetPelanggan)
		protected.GET("/pelanggan/:id", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.GetPelangganByID)
		protected.POST("/pelanggan", middleware.RequireRole(policy.Admin), controllers.CreatePelanggan)
		protected.PUT("/pelanggan/:id", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.UpdatePelanggan)
		protected.DELETE("/pelanggan/:id", middleware.RequireRole(policy.Admin), controllers.DeletePelanggan)

		// Route Order (Bisa diakses Pelanggan & Admin)
		protected.GET("/orders", middleware.RequireRole(policy.Admin, policy.Pelanggan), controllers.GetAllOrders)
		protected.POST("/orders", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.CreateOrder)
		protected.GET("/orders/pelanggan/:id_pelanggan", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.GetOrdersByPelanggan)
		protected.PUT("/orders/cancel/:id_order", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.CancelOrder)

		// Route khusus Update Status (HANYA Admin)
		protected.PUT("/orders/status/:id_order", middleware.RequireRole(policy.Admin), controllers.UpdateStatusOrder)

		// Route Pembayaran (Bisa diakses Pelanggan & Admin)
		// -----------------------------------------------------------------------
		// DITAMBAHKAN: GET /pembayaran untuk mengambil riwayat transaksi
		protected.GET("/pembayaran", middleware.RequireRole(policy.Admin, policy.Pelanggan), controllers.GetPembayaran)
		// -----------------------------------------------------------------------
		protected.POST("/pembayaran", middleware.RequireRole(policy.Pelanggan, policy.Admin), controllers.BayarPesanan)
	}

	r.Run(":8080")
}
