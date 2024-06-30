package controllers

import (
	"fmt"
	"golang-app/app/models"
	"golang-app/database"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (r *UserController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "user.html", data)
}

func (r *UserController) GetUsers(c *gin.Context) {
	var data []models.User
	database.DB.Find(&data)

	draw, _ := strconv.Atoi(c.Query("draw"))
	start, _ := strconv.Atoi(c.Query("start"))
	length, _ := strconv.Atoi(c.Query("length"))
	search := c.Query("search[value]")

	filteredData := make([]models.User, 0)
	for _, user := range data {
		v := reflect.ValueOf(user)
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i).Interface()
			if strings.Contains(strings.ToLower(fmt.Sprintf("%v", fieldValue)), strings.ToLower(search)) {
				filteredData = append(filteredData, user)
				break
			}
		}
	}

	totalRecords := len(filteredData)

	orderColumnIndex, _ := strconv.Atoi(c.Query("order[0][column]"))
	orderDirection := c.Query("order[0][dir]")

	switch orderColumnIndex {
	case 0:
		if orderDirection == "asc" {
			sort.Slice(filteredData, func(i, j int) bool {
				return filteredData[i].ID < filteredData[j].ID
			})
		} else {
			sort.Slice(filteredData, func(i, j int) bool {
				return filteredData[i].ID > filteredData[j].ID
			})
		}
	}

	end := start + length
	if end > totalRecords {
		end = totalRecords
	}
	pagedData := filteredData[start:end]

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    len(data),
		"recordsFiltered": totalRecords,
		"data":            pagedData,
	})
}

func (r *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	user.Password = "" // Hide password in response
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (r *UserController) InsertData(c *gin.Context) {
	var input struct {
		Name                 string `form:"name" json:"name" binding:"required"`
		Email                string `form:"email" json:"email" binding:"required,email"`
		NoTelepon            string `form:"no_telepon" json:"no_telepon" binding:"required"`
		Alamat               string `form:"alamat" json:"alamat" binding:"required"`
		Password             string `form:"password" json:"password" binding:"required,min=6"`
		PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind data", "message": err.Error()})
		return
	}

	if input.Password != input.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password confirmation does not match", "message": "Password confirmation does not match"})
		return
	}

	id := uuid.New().String()
	user := models.User{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		NoTelepon: input.NoTelepon,
		Alamat:    input.Alamat,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password", "message": err.Error()})
		return
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "for key 'users.uix_users_email'") {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Inserted"})
}

func (r *UserController) UpdateData(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name                 string `form:"name" json:"name" binding:"required"`
		Email                string `form:"email" json:"email" binding:"required,email"`
		NoTelepon            string `form:"no_telepon" json:"no_telepon" binding:"required"`
		Alamat               string `form:"alamat" json:"alamat" binding:"required"`
		Password             string `form:"password" json:"password"`
		PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	user.Email = input.Email
	user.Name = input.Name
	user.NoTelepon = input.NoTelepon
	user.Alamat = input.Alamat

	if input.Password != input.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password confirmation does not match", "message": "Password confirmation does not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Save(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "for key 'users.uix_users_email'") {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (r *UserController) DeleteData(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Deleted"})
}
