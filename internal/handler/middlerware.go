package handler

import (
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func AuthenticateMiddleware(authSvc auth.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// ✅ ตรวจ token
		if err := authSvc.VerifyToken(tokenStr); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// ✅ decode token อีกครั้งเพื่อดึง userID (แม้จะไม่ trust ข้อมูล 100%)
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(appcore_config.Config.SecretKey), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["sub"].(string)
			c.Set("user_id", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		c.Next()
	}
}



func AuthorizeMiddleware(authSvc auth.Jwt, allowedRoles ...string) gin.HandlerFunc {
	logrus.Info("in middleware")
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		userID := userIDVal.(string)

		// ✅ ดึง user จาก cache หรือ DB
		role, err := authSvc.GetRoleUserByID(c, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user not found"})
			return
		}

		// ✅ ตรวจว่า role ตรงไหม
		for _, allowed := range allowedRoles {
			if *role == allowed {
				logrus.Info("in middleware check role")
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
