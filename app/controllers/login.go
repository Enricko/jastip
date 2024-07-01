package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type LoginController struct {
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

func (r *LoginController) Index(c *gin.Context) {
	data := gin.H{
		"title":     "User Management",
		"csrfToken": csrf.GetToken(c),
	}
	c.HTML(http.StatusOK, "loginpage.html", data)
}
