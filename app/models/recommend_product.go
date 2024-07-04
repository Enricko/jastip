package models

import "time"

type StatusRecommend string

const (
	TidakAktif             StatusRecommend = "aktif"
	Aktif        StatusRecommend = "aktif"
)

type RecommendProduct struct {
	ID             string             `gorm:"primary_key;not null;" json:"id,omitempty"`
	IdBarang       string             `gorm:"not null;index" json:"id_barang" binding:"required"`
	Deskripsi string             `gorm:"TEXT;not null" json:"deskripsi" binding:"required"`
	Status         StatusProcessOrder `gorm:"type:enum('aktif','status');" json:"status" binding:"required"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`

	Barang Barang `gorm:"foreignkey:IdBarang;association_foreignkey:ID" json:"barang"`
}
