package middleware

import (
	"net/http"

	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.ValidateToken(c)
		if err != nil {
			utils.APIResponse(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", uint(claims["user_id"].(float64)))
			c.Set("userRole", claims["role"].(string))
		} else {
			utils.APIResponse(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			utils.APIResponse(c, http.StatusForbidden, "Role not found in token", nil)
			c.Abort()
			return
		}

		if role.(string) != requiredRole {
			utils.APIResponse(c, http.StatusForbidden, "Insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
