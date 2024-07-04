package controllers

import (
	"fmt"
	"golang-app/app/models"
	"golang-app/database"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"

)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type OrderController struct {
}

func NewOrderController() *OrderController {
	return &OrderController{}
}

func (r *OrderController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "Order Page",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "order.html", data)
}

func (r *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order

	if err := database.DB.Preload("User").First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (r *OrderController) GetOrders(c *gin.Context) {
	var data []models.Order

	result := database.DB.Preload("User").Preload("Barang").Find(&data)
	if result.Error != nil {
		// Log or handle the error
		fmt.Println("Error loading associated user data:", result.Error)
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	draw, _ := strconv.Atoi(c.Query("draw"))
	start, _ := strconv.Atoi(c.Query("start"))
	length, _ := strconv.Atoi(c.Query("length"))
	search := c.Query("search[value]")

	filteredData := make([]models.Order, 0)
	for _, order := range data {
		v := reflect.ValueOf(order)
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i).Interface()
			if strings.Contains(strings.ToLower(fmt.Sprintf("%v", fieldValue)), strings.ToLower(search)) {
				filteredData = append(filteredData, order)
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

func (r *OrderController) CreateOrder(c *gin.Context) {
	var input struct {
		NamaBarang string `form:"nama_barang" json:"nama_barang" binding:"required"`
		URL        string `form:"url" json:"url" binding:"required"`
		Alamat     string `form:"alamat" json:"alamat" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind data"})
		return
	}

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

	// Handle image upload
	file, err := c.FormFile("gambar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload image"})
		return
	}

	// Create directory if it doesn't exist
	dir := "public/upload/assets/image/order"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create directory"})
		return
	}

	// Save the file with a unique name
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filepath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
		return
	}

	id := uuid.New().String()
	idBarang := uuid.New().String()
	barang := models.Barang{
		ID:         idBarang,
		NamaBarang: input.NamaBarang,
		URL:        input.URL,
		Gambar:     filepath,
	}
	order := models.Order{
		ID:       id,
		IdUser:   user.ID,
		IdBarang: idBarang,
		Alamat:   input.Alamat,
		Status:   models.Status("baru"),
	}

	if err := database.DB.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

// UpdateOrder handles updating an existing order with image upload
func (r *OrderController) UpdateOrder(c *gin.Context) {
	var input struct {
		NamaBarang string `form:"nama_barang" json:"nama_barang" binding:"required"`
		URL        string `form:"url" json:"url" binding:"required"`
		Alamat     string `form:"alamat" json:"alamat" binding:"required"`
	}

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

	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("Barang").First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	if order.Status == models.Terima {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot update a confirmed order"})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind data"})
		return
	}

	// Handle image upload
	file, err := c.FormFile("gambar")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload image"})
		return
	}

	if file != nil {
		// Delete the old image file
		if order.Barang.Gambar != "" {
			if err := os.Remove(order.Barang.Gambar); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete old image"})
				return
			}
		}

		// Create directory if it doesn't exist
		dir := "public/upload/assets/image/order"
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create directory"})
			return
		}

		// Save the file with a unique name
		filename := uuid.New().String() + filepath.Ext(file.Filename)
		filepath := filepath.Join(dir, filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
			return
		}

		// Update the image path in the order
		order.Barang.Gambar = filepath
	}

	order.Barang.NamaBarang = input.NamaBarang
	order.Barang.URL = input.URL
	order.Alamat = input.Alamat

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder handles deleting an existing order
func (r *OrderController) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("Barang").First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	// Delete the image file
	if order.Barang.Gambar != "" {
		if err := os.Remove(order.Barang.Gambar); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete image"})
			return
		}
	}

	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// ConfirmOrder handles confirming an order
func (r *OrderController) ConfirmOrder(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		HargaBarang string `form:"harga_barang" json:"harga_barang" binding:"required"`
		Ongkir      string `form:"ongkir" json:"ongkir" binding:"required"`
	}

	var order models.Order
	if err := database.DB.First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	if order.Status == models.Terima {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot update a confirmed order"})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order confirmed successfully"})
}

// DeclineOrder handles declining an order
func (r *OrderController) DeclineOrder(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", models.Tolak).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order declined successfully"})
}
