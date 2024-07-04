package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type LuarNegeriController struct {
}

func NewLuarNegeriController() *LuarNegeriController {
	return &LuarNegeriController{}
}

func (r *LuarNegeriController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "luar_negeri.html", data)
}
