package middleware

import (
	"go-rebuild/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthenMiddleware is a middleware that checks for a valid JWT token
func JWTAuthenMiddleware(authSvc auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.GetHeader("Authorization")
		token := strings.TrimPrefix(s, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized empty token"})
			c.Abort()
			return
		}

		userID, err := authSvc.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}