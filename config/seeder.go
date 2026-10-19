package config

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	IDAdmin  uint   `gorm:"primaryKey;column:id_admin"`
	Nama     string `gorm:"column:nama"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

// TableName menentukan nama tabel di MariaDB
func (Admin) TableName() string {
	return "admin"
}

func SeedAdmin() {
	var admin Admin

	// 1. Ambil data kredensial dari environment variable (.env)
	email := os.Getenv("SEED_ADMIN_EMAIL")
	if email == "" {
		email = "admin@laundry.com" // Default fallback jika .env belum diisi
	}

	password := os.Getenv("SEED_ADMIN_PASSWORD")
	if password == "" {
		password = "admin123" // Default fallback
	}

	nama := os.Getenv("SEED_ADMIN_NAMA")
	if nama == "" {
		nama = "Admin Laundry" // Default fallback
	}

	// 2. Cari admin berdasarkan email dari .env
	err := DB.Where("email = ?", email).First(&admin).Error

	if err != nil {
		// Jika admin belum ada, buat akun baru
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			log.Println("Gagal meng-hash password admin:", err)
			return
		}

		admin = Admin{
			Nama:     nama,
			Email:    email,
			Password: string(hashedPassword),
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Gagal membuat admin:", err)
		} else {
			log.Println("Admin berhasil dibuat!")
		}

		return
	}

	// 3. Jika admin sudah ada, perbarui data & password sesuai .env
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Println("Gagal meng-hash password:", err)
		return
	}

	admin.Nama = nama
	admin.Password = string(hashedPassword)

	if err := DB.Save(&admin).Error; err != nil {
		log.Println("Gagal memperbarui password admin:", err)
	} else {
		log.Println("Password admin berhasil diperbarui!")
	}
}