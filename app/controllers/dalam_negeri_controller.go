package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type DalamNegeriController struct {
}

func NewDalamNegeriController() *DalamNegeriController {
	return &DalamNegeriController{}
}

func (r *DalamNegeriController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "dalam_negeri.html", data)
}
