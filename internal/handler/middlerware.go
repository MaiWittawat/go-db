package handler

import (
	"errors"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		if err := authSvc.VerifyToken(tokenStr); err != nil {
			log.WithError(err).WithFields(baseLogFields)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// decode token อีกครั้งเพื่อดึง userID 
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(appcore_config.Config.SecretKey), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["sub"].(string)
			c.Set("user_id", userID)
		} else {
			log.WithError(errors.New("failed to mapclaims")).WithFields(baseLogFields)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user not found e"})
			return
		}

		// ตรวจว่า role ตรงไหม
		for _, allowed := range allowedRoles {
			if *role == allowed {
				log.Info("in middleware check role pass")
				c.Next()
				return
			}
		}

		log.WithError(errors.New("role not match")).WithFields(baseLogFields)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no permissions"})
	}
}
