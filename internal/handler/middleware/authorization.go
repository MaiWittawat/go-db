package middleware

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func AuthorizeMiddleware(authSvc auth.AuthService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized empty user id"})
			c.Abort()
			return
		}

		var user model.User
		if err := authSvc.GetUserByID(c.Request.Context(), userID, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		log.Println("userID:", userID, "user.ID:", user.ID)
		if userID != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden id not match"})
			c.Abort()
			return
		}


		for _, role := range allowedRoles {
			log.Println("user role:", user.Role, "allowed role:", role)
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden permission denied"})
		c.Abort()
	}
}