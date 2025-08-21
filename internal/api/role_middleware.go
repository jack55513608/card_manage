package api

import (
	"card_manage/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware restricts access to handlers to specific roles.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the payload from the context
		payload, exists := c.Get(AuthorizationPayloadKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
			return
		}

		claims, ok := payload.(*service.CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid payload type"})
			return
		}

		// Check if the user's role is in the list of allowed roles
		isAllowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("role '%s' is not allowed to access this resource", claims.Role)})
			return
		}

		c.Next()
	}
}
