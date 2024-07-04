package models

import "time"

type StatusProcessOrder string

const (
	BelumBayar        StatusProcessOrder = "belum bayar"
	Lunas             StatusProcessOrder = "lunas"
	ProsesLuarNegeri  StatusProcessOrder = "proses luar negeri"
	ProsesDalamNegeri StatusProcessOrder = "proses dalam negeri"
	Selesai           StatusProcessOrder = "selesai"
)

type ProcessOrder struct {
	ID             string             `gorm:"primary_key;not null;" json:"id,omitempty"`
	IdUser         string             `gorm:"not null;index" json:"id_user" binding:"required"`
	IdBarang       string             `gorm:"not null;index" json:"id_barang" binding:"required"`
	PaymentGateway string             `gorm:"TEXT;not null" json:"payment_gateway" binding:"required"`
	HargaJasa      uint64             `gorm:"not null" json:"harga_jasa" binding:"required"`
	HargaOngkir    uint64             `gorm:"not null" json:"harga_ongkir" binding:"required"`
	NoResi         string             `gorm:"varchar(300)" json:"no_resi" binding:"required"`
	Status         StatusProcessOrder `gorm:"type:enum('belum bayar','lunas','proses luar negeri','proses dalam negeri','selesai');" json:"status" binding:"required"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`

	User   User   `gorm:"foreignkey:IdUser;association_foreignkey:ID" json:"user"`
	Barang Barang `gorm:"foreignkey:IdBarang;association_foreignkey:ID" json:"barang"`
}
