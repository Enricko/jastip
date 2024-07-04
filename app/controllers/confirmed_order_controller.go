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
	csrf "github.com/utrack/gin-csrf"
)

type ConfirmedOrderController struct {
}

func NewConfirmedOrderController() *ConfirmedOrderController {
	return &ConfirmedOrderController{}
}

func (r *ConfirmedOrderController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "Confirmed Order Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "confirmed_order.html", data)
}

func (r *ConfirmedOrderController) GetConfiredOrders(c *gin.Context) {
	var data []models.ProcessOrder

	result := database.DB.Preload("User").Preload("Barang").Where("status = ?", "belum bayar").Find(&data)
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

	filteredData := make([]models.ProcessOrder, 0)
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
