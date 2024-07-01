package api

import (
	"golang-app/app/models"
	"golang-app/database"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (r *AuthController) Register(c *gin.Context) {
	var input struct {
		Name                 string `form:"name" json:"name" binding:"required"`
		Email                string `form:"email" json:"email" binding:"required,email"`
		NoTelepon            string `form:"no_telepon" json:"no_telepon" binding:"required"`
		Alamat               string `form:"alamat" json:"alamat" binding:"required"`
		Password             string `form:"password" json:"password" binding:"required,min=6"`
		PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "err": err.Error()})
		return
	}

	if input.Password != input.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Passwords do not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	user := models.User{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Email:     input.Email,
		NoTelepon: input.NoTelepon,
		Alamat:    input.Alamat,
		Password:  string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "for key 'users.uix_users_email'") {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (r *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "err": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error finding user"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	expirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid token expiration configuration"})
		return
	}

	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)
	claims := &Claims{
		Email: input.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func (r *AuthController) Logout(c *gin.Context) {
	// Logout logic can be implemented if needed, such as invalidating the token on the server side.
	// In this simple implementation, we just return a success message.
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func (r *AuthController) ValidateToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, "Bearer "), claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token is valid", "user": user})
}
