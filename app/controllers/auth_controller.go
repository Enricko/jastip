package controllers

import (
	"golang-app/app/models"
	"golang-app/database"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (a *AuthController) LoginPage(c *gin.Context) {
	data := gin.H{
		"title":     "Login",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "loginpage.html", data)
}

func (ac *AuthController) Logout(c *gin.Context) {
	// Invalidate the JWT token by setting an expired cookie
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)

	c.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginData struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Administrator
	if err := database.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID, string(user.Level))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token in a cookie
	c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true) // Expires in 24 hours

	c.Redirect(http.StatusFound, "/dashboard")
}

func generateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": userID,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret")) // Replace "secret" with your actual secret key
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
