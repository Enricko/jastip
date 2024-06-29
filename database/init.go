package database

import (
    "os"
	"fmt"
	"golang-app/app/models"

    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // MYSQL

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

	DB.AutoMigrate(&models.Adminstrator{})
}