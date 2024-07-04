package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type BaruController struct {
}

func NewBaruController() *BaruController {
	return &BaruController{}
}

func (r *BaruController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "baru.html", data)
}
