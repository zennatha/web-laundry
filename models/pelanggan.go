package models

import "time"

type Pelanggan struct {
	IDPelanggan uint      `gorm:"primaryKey;column:id_pelanggan" json:"id_pelanggan"`
	Nama        string    `gorm:"column:nama" json:"nama"`
	NoHp        string    `gorm:"column:no_hp" json:"no_hp"`
	Alamat      string    `gorm:"column:alamat" json:"alamat"`
	Email       string    `gorm:"column:email;unique" json:"email"`
	Password    string    `gorm:"column:password" json:"password,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Pelanggan) TableName() string {
	return "pelanggan"
}
