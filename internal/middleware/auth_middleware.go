package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/enum"
	auth "github.com/ners1us/order-service/internal/service"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrNoAuthToken.Error()})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrWrongTokenFormat.Error()})
			c.Abort()
			return
		}
		tokenStr := parts[1]
		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrInvalidToken.Error()})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
