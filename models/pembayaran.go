package models

import "time"

type Pembayaran struct {
	IDPembayaran     uint      `gorm:"primaryKey;column:id_pembayaran" json:"id_pembayaran"`
	IDOrder          uint      `gorm:"column:id_order" json:"id_order"`
	IDPelanggan      uint      `gorm:"column:id_pelanggan" json:"id_pelanggan"`
	MetodePembayaran string    `gorm:"column:metode_pembayaran" json:"metode_pembayaran"`
	TotalHarga       int       `gorm:"column:total_harga" json:"total_harga"`
	Status           string    `gorm:"column:status;default:Pending" json:"status"`
	Tanggal          time.Time `gorm:"column:tanggal;autoCreateTime" json:"tanggal"`
}

func (Pembayaran) TableName() string {
	return "pembayaran"
}