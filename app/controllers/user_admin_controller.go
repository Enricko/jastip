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

type UserAdminController struct {
}

func NewUserAdminController() *UserAdminController {
	return &UserAdminController{}
}

func (r *UserAdminController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "Admin User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "user_admin.html", data)
}

func (r *UserAdminController) GetAdminUsers(c *gin.Context) {
	var data []models.Adminstrator
	database.DB.Find(&data)

	draw, _ := strconv.Atoi(c.Query("draw"))
	start, _ := strconv.Atoi(c.Query("start"))
	length, _ := strconv.Atoi(c.Query("length"))
	search := c.Query("search[value]")

	filteredData := make([]models.Adminstrator, 0)
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

func (r *UserAdminController) GetAdminUser(c *gin.Context) {
	id := c.Param("id")

	var user models.Adminstrator
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Admin User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (r *UserAdminController) InsertData(c *gin.Context) {
	var input struct {
		Name                 string `form:"name" json:"name" binding:"required"`
		Email                string `form:"email" json:"email" binding:"required,email"`
		NoTelepon            string `form:"no_telepon" json:"no_telepon" binding:"required"`
		Level                string `form:"level" json:"level" binding:"required"`
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

	var existingAdmin models.Adminstrator
	if err := database.DB.Where("email = ?", input.Email).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists", "message": "Email already exists"})
		return
	}

	// Generate a random UUID for ID
	id := uuid.New().String()
	admin := models.Adminstrator{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		NoTelepon: input.NoTelepon,
		Level:     models.Level(input.Level),
	}

	// Encrypt password before inserting into database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password", "message": err.Error()})
		return
	}
	admin.Password = string(hashedPassword)

	// Insert the data into the database
	if err := database.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Inserted"})
}

func (r *UserAdminController) UpdateData(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name                 string `form:"name" json:"name" binding:"required"`
		Email                string `form:"email" json:"email" binding:"required,email"`
		NoTelepon            string `form:"no_telepon" json:"no_telepon" binding:"required"`
		Level                string `form:"level" json:"level" binding:"required"`
		Password             string `form:"password" json:"password" binding:"required,min=6"`
		PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var admin models.Adminstrator
	if err := database.DB.First(&admin, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Admin user not found"})
		return
	}

	if input.Password != input.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password confirmation does not match", "message": "Password confirmation does not match"})
		return
	}

	if admin.Email != input.Email {
		var existingAdmin models.Adminstrator
		if err := database.DB.Where("email = ?", input.Email).First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists", "message": "Email already exists"})
			return
		}
	}

	admin.Email = input.Email
	admin.Name = input.Name
	admin.NoTelepon = input.NoTelepon
	admin.Level = models.Level(input.Level)

	// Hash the updated password before saving
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	admin.Password = hashedPassword

	if err := database.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update admin user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin user updated successfully"})
}

func (r *UserAdminController) DeleteData(c *gin.Context) {
	id := c.Param("id")

	var user models.Adminstrator
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Admin User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin user Deleted"})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
