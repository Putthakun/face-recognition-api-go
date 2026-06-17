package middleware

import (
	"net/http"
	"strings"

	"github.com/Putthakun/face-recognition-api-go/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyEmpID = "empId"
	ContextKeyRole  = "role"
)

func Auth(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtService.Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		c.Set(ContextKeyEmpID, claims.EmpID)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		role, _ := c.Get(ContextKeyRole)
		if !allowed[role.(string)] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "insufficient permissions"})
			return
		}
		c.Next()
	}
}
