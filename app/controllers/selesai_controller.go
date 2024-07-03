package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type SelesaiController struct {
}

func NewSelesaiController() *SelesaiController {
	return &SelesaiController{}
}

func (r *SelesaiController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "selesai.html", data)
}
