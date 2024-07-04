package main

import (
	"golang-app/app/middleware"
	"golang-app/database"
	"golang-app/routes"
	"log"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.Default())

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

	r.Run("0.0.0.0:8080")

	defer database.DB.Close()
}
