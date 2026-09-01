package config

import (
	"log"

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
	var count int64
	DB.Model(&Admin{}).Where("email = ?", "zenncatt1@gmail.com").Count(&count)

	// Jika admin belum ada, buat admin baru dengan password ter-hash
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Gagal meng-hash password admin:", err)
			return
		}

		admin := Admin{
			Nama:     "zen",                 // <-- Nama diubah menjadi zen
			Email:    "zenncatt1@gmail.com", // <-- Email diubah
			Password: string(hashedPassword),
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Gagal membuat seed admin:", err)
		} else {
			log.Println("Seeder Admin berhasil dibuat!")
		}
	}
}
