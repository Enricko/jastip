package models

import "time"

type Status string

const (
	Terima Status = "terima"
	Tolak  Status = "tolak"
	Baru  Status = "baru"
)

type Order struct {
	ID         string    `gorm:"primary_key;not null;" json:"id,omitempty"`
	IdUser     string    `gorm:"not null;index" json:"id_user" binding:"required"`
	NamaBarang string    `gorm:"varchar(300);not null;" json:"nama_barang" binding:"required"`
	URL        string    `gorm:"TEXT;not null;" json:"url" binding:"required"`
	Gambar     string    `gorm:"TEXT;not null;" binding:"required" json:"gambar"`
	Alamat     string    `gorm:"TEXT;" json:"alamat" binding:"required"`
	Status     Status    `gorm:"type:enum('terima','tolak','baru');" json:"status" binding:"required`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User User `gorm:"foreignkey:IdUser;association_foreignkey:ID" json:"user"`
}
