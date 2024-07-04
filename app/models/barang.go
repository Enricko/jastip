package models

import "time"

type Barang struct {
	ID          string    `gorm:"primary_key;not null;" json:"id,omitempty"`
	NamaBarang  string    `gorm:"varchar(300);not null;" json:"nama_barang" binding:"required"`
	HargaBarang uint64    `json:"harga_barang"`
	URL         string    `gorm:"TEXT;not null;" json:"url" binding:"required"`
	Berat         string    `gorm:"TEXT;not null;" json:"berat"`
	Gambar      string    `gorm:"TEXT;not null;" binding:"required" json:"gambar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
