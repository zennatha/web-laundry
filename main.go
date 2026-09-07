package main

import (
	"laundry-app/config"
	"laundry-app/controllers"
	"laundry-app/middleware"
	"laundry-app/policy"
	"time"

	"github.com/gin-contrib/cors" // Import middleware CORS
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Hubungkan ke database MariaDB
	config.ConnectDatabase()

	// 2. Jalankan seeder admin
	config.SeedAdmin()

	r := gin.Default()

	// 3. Konfigurasi CORS Lengkap
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // URL Frontend
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
	r.POST("/api/register/admin", controllers.RegisterAdmin) // <-- ROUTE DITAMBAHKAN DI SINI
	r.POST("/api/login", controllers.Login)
	r.POST("/api/login/admin", controllers.LoginAdmin)
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
