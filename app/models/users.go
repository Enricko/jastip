package models

import "time"

type User struct {
	ID        string    `gorm:"primary_key;not null;" json:"id,omitempty"`
	Name      string    `gorm:"varchar(300);not null;" json:"name" binding:"required"`
	Email     string    `gorm:"unique_index;not null;" binding:"required" json:"email"`
	NoTelepon string    `gorm:"varchar(300);not null;" json:"no_telepon" binding:"required"`
	Alamat    string    `gorm:"varchar(300);not null;" json:"alamat" binding:"required"`
	Password  string    `gorm:"varchar(300);not null;" json:"password" binding:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
