package middleware

import (
	"net/http"
	// "strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"message": "You Dont have permission to access this resource"})
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("admin_id", claims["admin_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

// func ApiAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
// 			c.Abort()
// 			return
// 		}

// 		// tokenString := strings.Split(authHeader, "Bearer ")[1]
// 		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("secret"), nil
// 		})

// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			userID := claims["id"].(string)
// 			c.Set("userID", userID)
// 		} else {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

func RoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != role {
			c.JSON(http.StatusForbidden, gin.H{"message": "You don't have permission to access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RedirectIfLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err == nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})

			if err == nil && token.Valid {
				c.Redirect(http.StatusFound, "/dashboard")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
