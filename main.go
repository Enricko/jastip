package main

import (
	"golang-app/database"
	"golang-app/routes"
	"golang-app/app/middleware"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Init()

	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	excludedPaths := []string{
		"/api/auth",
		// Add more paths as needed
	}

	// Use the NoCSRF middleware with the exclusion list
	r.Use(middleware.NoCSRF(excludedPaths))

	r.LoadHTMLGlob("templates/**/*")
	routes.SetupRouter(r)

	r.Static("/public", "./public")

	r.Run("127.0.0.1:8080")

	defer database.DB.Close()
}
