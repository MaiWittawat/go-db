package handler

import (
	"errors"
	"go-rebuild/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AuthenticateMiddleware(authSvc auth.Jwt) gin.HandlerFunc {
	var baseLogFields = log.Fields{
		"layer":     "middleware",
		"operation": "authenticate_middleware",
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.WithError(errors.New("no string bearer")).WithFields(baseLogFields)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจ token
		claims, err := authSvc.VerifyToken(tokenStr)
		if err != nil {
			log.WithError(err).WithFields(baseLogFields)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims == nil || claims.Subject == "" {
			log.WithError(errors.New("failed to mapclaims")).WithFields(baseLogFields)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		
		userID := claims.Subject
		c.Set("user_id", userID)
		c.Next()
	}
}

func AuthorizeMiddleware(authSvc auth.Jwt, allowedRoles ...string) gin.HandlerFunc {
	var baseLogFields = log.Fields{
		"layer":     "middleware",
		"operation": "authorize_middleware",
	}
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		userID := userIDVal.(string)

		// ดึง user จาก cache หรือ DB
		role, err := authSvc.GetRoleUserByID(c, userID)
		if err != nil {
			log.WithError(err).WithFields(baseLogFields).Error("user not found")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user not found"})
			return
		}

		// ตรวจว่า role ตรงไหม
		for _, allowed := range allowedRoles {
			if *role == allowed {
				log.Info("[Middleware]: Role pass")
				c.Next()
				return
			}
		}

		log.WithError(errors.New("role not match")).WithFields(baseLogFields)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no permissions"})
	}
}
