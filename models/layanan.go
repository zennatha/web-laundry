package models

type Layanan struct {
	IDLayanan   uint   `gorm:"primaryKey;column:id_layanan" json:"id_layanan"`
	NamaLayanan string `gorm:"column:nama_layanan" json:"nama_layanan"`
	Harga       int    `gorm:"column:harga" json:"harga"`
	Estimasi    string `gorm:"column:estimasi" json:"estimasi"`
}

func (Layanan) TableName() string {
	return "layanan"
}
