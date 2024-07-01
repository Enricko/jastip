package database

import (
	"golang-app/app/models"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

)

// SeedAdminstrators seeds initial data for Adminstrator model
func SeedAdminstrators() {
	// Define initial data
	admins := []models.Administrator{
		{
			ID:        uuid.New().String(),
			Name:      "Admin One",
			Email:     "admin1@gmail.com",
			NoTelepon: "1234567890",
			Level:     models.Admin,
			Password:  hashPassword("admin"),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Owner One",
			Email:     "owner1@gmail.com",
			NoTelepon: "0987654321",
			Level:     models.Owner,
			Password:  hashPassword("admin"),
		},
	}

	for _, admin := range admins {
		err := DB.Where(models.Administrator{Email: admin.Email}).FirstOrCreate(&admin).Error
		if err != nil {
			log.Printf("Cannot seed data: %v", err)
		}
	}
}

// hashPassword hashes a password string
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}
