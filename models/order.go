package models

type Order struct {
	IDOrder     uint    `gorm:"primaryKey;column:id_order" json:"id_order"`
	IDPelanggan uint    `gorm:"column:id_pelanggan" json:"id_pelanggan"`
	IDAdmin     *uint   `gorm:"column:id_admin" json:"id_admin"`
	Layanan     string  `gorm:"column:layanan" json:"layanan"`
	Berat       float64 `gorm:"column:berat" json:"berat"`
	Harga       int     `gorm:"column:harga" json:"harga"`
	Status      string  `gorm:"column:status;default:Proses" json:"status"`
}

func (Order) TableName() string {
	return "order"
}
