package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type ProdukRekomendasiController struct {
}

func NewProdukRekomendasiController() *ProdukRekomendasiController {
	return &ProdukRekomendasiController{}
}

func (r *ProdukRekomendasiController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "Admin User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "produk_rekomendasi.html", data)
}
