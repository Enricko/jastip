package database

import (
	"fmt"
	"golang-app/app/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // MYSQL
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

var DB *gorm.DB

func Init() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	// dbSSLMode := os.Getenv("DB_SSL_MODE")

	dbConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	var errDb error
	// Connect to PostgreSQL database
	DB, errDb = gorm.Open("mysql", dbConnection)
	if errDb != nil {
		panic(errDb.Error())
	}

	err := DB.AutoMigrate(&models.Administrator{}, &models.User{}, &models.Barang{}, &models.Order{},&models.ProcessOrder{},&models.RecommendProduct{}).Error
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := DB.Model(&models.Order{}).AddForeignKey("id_barang", "barangs(id)", "CASCADE", "CASCADE").AddForeignKey("id_user", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.Fatal("Failed to set up foreign key:", err)
	}
	if err := DB.Model(&models.ProcessOrder{}).AddForeignKey("id_barang", "barangs(id)", "CASCADE", "CASCADE").AddForeignKey("id_user", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.Fatal("Failed to set up foreign key:", err)
	}
	if err := DB.Model(&models.RecommendProduct{}).AddForeignKey("id_barang", "barangs(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.Fatal("Failed to set up foreign key:", err)
	}

	SeedAdminstrators()
}
