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
	var admin Admin

	// Cari admin berdasarkan email
	err := DB.Where("email = ?", "zenncatt1@gmail.com").First(&admin).Error

	if err != nil {
		// Jika admin belum ada, buat akun baru
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte("admin123"),
			bcrypt.DefaultCost,
		)

		if err != nil {
			log.Println("Gagal meng-hash password admin:", err)
			return
		}

		admin = Admin{
			Nama:     "zen",
			Email:    "zenncatt1@gmail.com",
			Password: string(hashedPassword),
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Gagal membuat admin:", err)
		} else {
			log.Println("Admin berhasil dibuat!")
		}

		return
	}

	// Jika admin sudah ada, reset password menjadi admin123
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Println("Gagal meng-hash password:", err)
		return
	}

	admin.Password = string(hashedPassword)

	if err := DB.Save(&admin).Error; err != nil {
		log.Println("Gagal memperbarui password admin:", err)
	} else {
		log.Println("Password admin berhasil diperbarui!")
	}
}
