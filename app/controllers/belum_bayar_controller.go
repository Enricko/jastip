package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type BelumBayarController struct {
}

func NewBelumBayarController() *BelumBayarController {
	return &BelumBayarController{}
}

func (r *BelumBayarController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "belum_bayar.html", data)
}
