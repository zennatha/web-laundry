package main

import (
	"laundry-app/config"
	"laundry-app/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	// Route Auth
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// Route Layanan
	r.GET("/api/layanan", controllers.GetLayanan)

	// Route Order / Pesanan
	r.POST("/api/orders", controllers.CreateOrder)
	r.GET("/api/orders/pelanggan/:id_pelanggan", controllers.GetOrdersByPelanggan)
	r.PUT("/api/orders/status/:id_order", controllers.UpdateStatusOrder)

	// Route Pembayaran (DITAMBAHKAN DI SINI)
	r.POST("/api/pembayaran", controllers.BayarPesanan)

	r.Run(":8080")
}
